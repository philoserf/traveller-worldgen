// Command worldgen generates Traveller worlds. Each Traveller edition is a
// subcommand — e.g. "worldgen classic -seed 42". Run "worldgen" with no
// arguments to list the available editions.
package main

import (
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

// editions maps a subcommand name to its runner. Adding an edition is one entry
// here plus its runner file (e.g. classic.go).
var editions = map[string]func(args []string, stdout, stderr io.Writer) int{
	"classic": runClassic,
	"mega":    runMega,
	"tne":     runTne,
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

// renderWorldsJSON marshals one world as an object, or several as an array, with
// two-space indentation and a trailing newline. Shared by every edition's runner
// so the JSON shape stays identical across editions.
func renderWorldsJSON[T any](worlds []T) (string, error) {
	var v any = worlds
	if len(worlds) == 1 {
		v = worlds[0]
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding json: %w", err)
	}
	return string(data) + "\n", nil
}

// renderWorlds turns the generated worlds into the requested output format,
// using the edition's per-world formatters for the uwp and text layouts. The
// json shape comes from each World's MarshalJSON via the shared renderWorldsJSON.
func renderWorlds[T any](format string, worlds []T, uwpLine, textBlock func(T) string) (string, error) {
	switch format {
	case "json":
		return renderWorldsJSON(worlds)
	case "uwp":
		var b strings.Builder
		for _, w := range worlds {
			b.WriteString(uwpLine(w))
			b.WriteByte('\n')
		}
		return b.String(), nil
	default: // text
		var b strings.Builder
		for i, w := range worlds {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(textBlock(w))
		}
		return b.String(), nil
	}
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
	if *n < 1 {
		errf(stderr, "worldgen: -n must be >= 1, got %d\n", *n)
		return 2
	}
	switch *format {
	case "text", "uwp", "json":
	default:
		errf(stderr, "worldgen: unknown -format %q (want text, uwp, or json)\n", *format)
		return 2
	}
	generate, err := validate()
	if err != nil {
		errf(stderr, "worldgen: %v\n", err)
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

	// One roller drives all N worlds, so a seed reproduces the whole batch.
	roller := dice.NewSeeded(*seed)
	worlds := make([]T, *n)
	for i := range worlds {
		worlds[i] = generate(roller)
	}

	out, err := renderWorlds(*format, worlds, uwpLine, textBlock)
	if err != nil {
		errf(stderr, "worldgen: %v\n", err)
		return 1
	}
	if _, err := io.WriteString(stdout, out); err != nil {
		errf(stderr, "worldgen: writing output: %v\n", err)
		return 1
	}
	return 0
}
