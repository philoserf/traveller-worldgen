// Command worldgen generates Traveller worlds. Each Traveller edition is a
// subcommand — e.g. "worldgen classic -seed 42". Run "worldgen" with no
// arguments to list the available editions.
package main

import (
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
