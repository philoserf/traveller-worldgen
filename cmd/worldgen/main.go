// Command worldgen generates Traveller worlds. Each Traveller edition is a
// subcommand — e.g. "worldgen classic -seed 42". Run "worldgen" with no
// arguments to list the available editions.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/philoserf/traveller-worldgen/dice"
)

// maxWorlds bounds the -n flag so an accidental huge value can't spin the CLI
// indefinitely. Worlds stream out one at a time (see streamWorlds), so a large
// batch no longer buffers in memory; the cap only limits total work. A full
// sector is 640 hexes, so this is generous headroom for any realistic batch.
const maxWorlds = 1_000_000

// editions maps a subcommand name to its runner. Adding an edition is one entry
// here plus its runner file (e.g. classic.go).
var editions = map[string]func(args []string, stdout, stderr io.Writer) int{
	"classic": runClassic,
	"mega":    runMega,
	"tne":     runTne,
	"t5":      runT5,
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run dispatches on the edition subcommand. It returns a process exit code:
// 0 on success, 1 on an output error, 2 on a usage error.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	switch args[0] {
	case "-h", "--help", "help":
		usage(stdout)
		return 0
	}
	runner, ok := editions[args[0]]
	if !ok {
		errf(stderr, "worldgen: unknown edition %q\n\n", args[0])
		usage(stderr)
		return 2
	}
	return runner(args[1:], stdout, stderr)
}

// usage lists the available editions.
func usage(w io.Writer) {
	names := slices.Sorted(maps.Keys(editions))
	errf(w, "usage: worldgen <edition> [flags]\n\neditions:\n  %s\n\nrun 'worldgen <edition> -h' for that edition's flags\n",
		strings.Join(names, "\n  "))
}

// errf writes a diagnostic to w, deliberately ignoring the write error (there
// is nowhere better to report a failure to write an error message).
func errf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}

// streamWorlds generates the n worlds one at a time and writes each to w in the
// requested format. Generating and writing per world keeps memory flat in n
// instead of materializing the whole batch and its rendered output up front — a
// large -n would otherwise buffer gigabytes. It routes through a bufio.Writer
// whose final Flush surfaces any write error; bufio records write errors stickily,
// so the intermediate writes deliberately ignore their returns. Per-world layout
// comes from the edition's uwp/text formatters; the json shape comes from each
// World's MarshalJSON.
func streamWorlds[T any](
	w io.Writer, format string, n int,
	roller dice.Roller, generate func(dice.Roller) T,
	uwpLine, textBlock func(T) string,
) error {
	bw := bufio.NewWriter(w)
	switch format {
	case "json":
		if err := streamWorldsJSON(bw, n, roller, generate); err != nil {
			return err
		}
	case "uwp":
		for range n {
			_, _ = bw.WriteString(uwpLine(generate(roller)))
			_ = bw.WriteByte('\n')
		}
	default: // text
		for i := range n {
			if i > 0 {
				_ = bw.WriteByte('\n')
			}
			_, _ = bw.WriteString(textBlock(generate(roller)))
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}

// streamWorldsJSON writes the worlds as JSON: a single object when n == 1, else an
// array. Unwrapping the lone n == 1 world to a bare object is a deliberate
// ergonomic choice for the common single-world invocation (consumers of a batch
// already handle the array); do not "fix" it into an always-array. The array
// framing is written by hand so each world is marshaled on its own, but the bytes
// match json.MarshalIndent over the whole slice (two-space indent, trailing
// newline) — an array element sits at indent depth 1, which MarshalIndent with
// prefix "  " reproduces.
func streamWorldsJSON[T any](bw *bufio.Writer, n int, roller dice.Roller, generate func(dice.Roller) T) error {
	if n == 1 {
		data, err := json.MarshalIndent(generate(roller), "", "  ")
		if err != nil {
			return fmt.Errorf("encoding json: %w", err)
		}
		_, _ = bw.Write(data)
		_ = bw.WriteByte('\n')
		return nil
	}
	_, _ = bw.WriteString("[\n")
	for i := range n {
		if i > 0 {
			_, _ = bw.WriteString(",\n")
		}
		data, err := json.MarshalIndent(generate(roller), "  ", "  ")
		if err != nil {
			return fmt.Errorf("encoding json: %w", err)
		}
		_, _ = bw.WriteString("  ")
		_, _ = bw.Write(data)
	}
	_, _ = bw.WriteString("\n]\n")
	return nil
}

// runEdition runs one edition's CLI. It owns the flags every edition shares
// (-seed, -format, -n) plus the parse/validate/seed/generate/render/write flow,
// leaving each edition to register any extra flags and supply its per-world
// formatters. bindGenerator registers extra flags on fs (before parsing) and
// returns a function, called after a successful parse, that validates those
// flags and returns the per-world generator — or an error (reported as a usage
// error) to abort. It returns a process exit code: 0 on success, 1 on an output
// error, 2 on a usage error.
func runEdition[T any](
	name string,
	args []string, stdout, stderr io.Writer,
	bindGenerator func(fs *flag.FlagSet) func() (func(dice.Roller) T, error),
	uwpLine, textBlock func(T) string,
) int {
	fs := flag.NewFlagSet("worldgen "+name, flag.ContinueOnError)
	fs.SetOutput(io.Discard) // suppress flag's auto-print; we route help/errors ourselves

	seed := fs.Int64("seed", 0, "random seed (default: time-based)")
	format := fs.String("format", "text", "output format: text, uwp, or json")
	n := fs.Int("n", 1, "number of independent worlds to generate")
	validate := bindGenerator(fs)

	if err := fs.Parse(args); err != nil {
		// -h/-help is a help request, not a usage error: print to stdout, exit 0
		// (matching the top-level `worldgen --help`).
		if errors.Is(err, flag.ErrHelp) {
			fs.SetOutput(stdout)
			fs.Usage()
			return 0
		}
		errf(stderr, "worldgen %s: %v\n", name, err)
		return 2
	}
	if *n < 1 || *n > maxWorlds {
		errf(stderr, "worldgen %s: -n must be between 1 and %d, got %d\n", name, maxWorlds, *n)
		return 2
	}
	switch *format {
	case "text", "uwp", "json":
	default:
		errf(stderr, "worldgen %s: unknown -format %q (want text, uwp, or json)\n", name, *format)
		return 2
	}
	generate, err := validate()
	if err != nil {
		errf(stderr, "worldgen %s: %v\n", name, err)
		return 2
	}

	// An explicit -seed (even -seed 0) is honored; otherwise seed from the clock.
	seeded := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			seeded = true
		}
	})
	if !seeded {
		*seed = time.Now().UnixNano()
	}

	// One roller drives all N worlds, so a seed reproduces the whole batch. Worlds
	// are generated and written one at a time, so a large -n streams to stdout in
	// flat memory rather than buffering the batch and its rendered output.
	roller := dice.NewSeeded(*seed)
	if err := streamWorlds(stdout, *format, *n, roller, generate, uwpLine, textBlock); err != nil {
		errf(stderr, "worldgen: %v\n", err)
		return 1
	}
	return 0
}
