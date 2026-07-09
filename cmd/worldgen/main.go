// Command worldgen generates Traveller worlds. Each Traveller edition is a
// subcommand — e.g. "worldgen classic -seed 42". Run "worldgen" with no
// arguments to list the available editions.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
)

// editions maps a subcommand name to its runner. Adding an edition is one entry
// here plus its runner file (e.g. classic.go).
var editions = map[string]func(args []string, stdout, stderr io.Writer) int{
	"classic": runClassic,
	"mega":    runMega,
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
	runEdition, ok := editions[args[0]]
	if !ok {
		errf(stderr, "worldgen: unknown edition %q\n\n", args[0])
		usage(stderr)
		return 2
	}
	return runEdition(args[1:], stdout, stderr)
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
	var (
		data []byte
		err  error
	)
	if len(worlds) == 1 {
		data, err = json.MarshalIndent(worlds[0], "", "  ")
	} else {
		data, err = json.MarshalIndent(worlds, "", "  ")
	}
	if err != nil {
		return "", fmt.Errorf("encoding json: %w", err)
	}
	return string(data) + "\n", nil
}
