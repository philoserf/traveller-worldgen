package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files in testdata/")

func runCapture(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errBuf bytes.Buffer
	code := run(args, &out, &errBuf)
	return code, out.String(), errBuf.String()
}

func TestClassicGoldenFormats(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		golden string
	}{
		{"text", []string{"classic", "-seed", "42"}, "seed42.text"},
		{"uwp", []string{"classic", "-seed", "42", "-format", "uwp"}, "seed42.uwp"},
		{"json", []string{"classic", "-seed", "42", "-format", "json"}, "seed42.json"},
		{"batch", []string{"classic", "-seed", "7", "-n", "3", "-format", "uwp"}, "seed7n3.uwp"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, out, errStr := runCapture(t, c.args...)
			if code != 0 {
				t.Fatalf("exit %d, stderr: %s", code, errStr)
			}
			path := filepath.Join("testdata", c.golden)
			if *update {
				if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden (run with -update to create): %v", err)
			}
			if out != string(want) {
				t.Errorf("output mismatch for %s:\n got:\n%s\nwant:\n%s", c.golden, out, want)
			}
		})
	}
}

func TestMegaGoldenFormats(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		golden string
	}{
		{"text", []string{"mega", "-seed", "42"}, "mega-seed42.text"},
		{"uwp", []string{"mega", "-seed", "42", "-format", "uwp"}, "mega-seed42.uwp"},
		{"json", []string{"mega", "-seed", "42", "-format", "json"}, "mega-seed42.json"},
		{"backwater", []string{"mega", "-seed", "5", "-n", "8", "-nature", "backwater", "-format", "uwp"}, "mega-seed5n8-backwater.uwp"},
		{"batch", []string{"mega", "-seed", "7", "-n", "3", "-format", "uwp"}, "mega-seed7n3.uwp"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, out, errStr := runCapture(t, c.args...)
			if code != 0 {
				t.Fatalf("exit %d, stderr: %s", code, errStr)
			}
			path := filepath.Join("testdata", c.golden)
			if *update {
				if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden (run with -update to create): %v", err)
			}
			if out != string(want) {
				t.Errorf("output mismatch for %s:\n got:\n%s\nwant:\n%s", c.golden, out, want)
			}
		})
	}
}

func TestMegaDeterministic(t *testing.T) {
	_, a, _ := runCapture(t, "mega", "-seed", "123")
	_, b, _ := runCapture(t, "mega", "-seed", "123")
	if a != b {
		t.Fatalf("same seed produced different output:\n%s\n---\n%s", a, b)
	}
}

func TestMegaJSONFields(t *testing.T) {
	code, out, errStr := runCapture(t, "mega", "-seed", "42", "-format", "json")
	if code != 0 {
		t.Fatalf("exit %d, stderr: %s", code, errStr)
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	for _, field := range []string{"uwp", "militaryBase", "gasGiants", "tradeCodes", "baseCode"} {
		if _, ok := obj[field]; !ok {
			t.Errorf("JSON missing %q field", field)
		}
	}
}

func TestMegaErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"unknown flag", []string{"mega", "-nope"}},
		{"bad format", []string{"mega", "-format", "xml"}},
		{"bad nature", []string{"mega", "-nature", "bogus"}},
		{"bad n", []string{"mega", "-n", "0"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, _, errStr := runCapture(t, c.args...)
			if code != 2 {
				t.Errorf("exit = %d, want 2", code)
			}
			if strings.TrimSpace(errStr) == "" {
				t.Error("expected a message on stderr")
			}
		})
	}
}

func TestMegaHelp(t *testing.T) {
	code, out, errStr := runCapture(t, "mega", "-h")
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	if !strings.Contains(out, "-nature") {
		t.Errorf("help should list the -nature flag on stdout, got: %s", out)
	}
	if strings.TrimSpace(errStr) != "" {
		t.Errorf("help should not write to stderr, got: %s", errStr)
	}
}

func TestClassicDeterministic(t *testing.T) {
	_, a, _ := runCapture(t, "classic", "-seed", "123")
	_, b, _ := runCapture(t, "classic", "-seed", "123")
	if a != b {
		t.Fatalf("same seed produced different output:\n%s\n---\n%s", a, b)
	}
}

func TestClassicJSONValid(t *testing.T) {
	code, out, errStr := runCapture(t, "classic", "-seed", "42", "-format", "json")
	if code != 0 {
		t.Fatalf("exit %d, stderr: %s", code, errStr)
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if obj["uwp"] == "" {
		t.Error("JSON missing uwp field")
	}
}

func TestClassicJSONArrayForBatch(t *testing.T) {
	code, out, _ := runCapture(t, "classic", "-seed", "42", "-n", "2", "-format", "json")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	var arr []map[string]any
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("expected JSON array for -n 2: %v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("array length = %d, want 2", len(arr))
	}
}

func TestClassicErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"unknown flag", []string{"classic", "-nope"}},
		{"bad format", []string{"classic", "-format", "xml"}},
		{"bad n", []string{"classic", "-n", "0"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, _, errStr := runCapture(t, c.args...)
			if code != 2 {
				t.Errorf("exit = %d, want 2", code)
			}
			if strings.TrimSpace(errStr) == "" {
				t.Error("expected a message on stderr")
			}
		})
	}
}

func TestClassicHelp(t *testing.T) {
	// `worldgen classic -h` is a help request, not an error: exit 0, flags to
	// stdout, nothing on stderr — consistent with the top-level `--help`.
	code, out, errStr := runCapture(t, "classic", "-h")
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	if !strings.Contains(out, "-seed") {
		t.Errorf("help should list the -seed flag on stdout, got: %s", out)
	}
	if strings.TrimSpace(errStr) != "" {
		t.Errorf("help should not write to stderr, got: %s", errStr)
	}
}

func TestDispatch(t *testing.T) {
	t.Run("no args lists editions", func(t *testing.T) {
		code, _, errStr := runCapture(t)
		if code != 2 {
			t.Errorf("exit = %d, want 2", code)
		}
		if !strings.Contains(errStr, "classic") {
			t.Errorf("usage should list editions, got: %s", errStr)
		}
	})
	t.Run("help to stdout", func(t *testing.T) {
		code, out, _ := runCapture(t, "--help")
		if code != 0 {
			t.Errorf("exit = %d, want 0", code)
		}
		if !strings.Contains(out, "classic") {
			t.Errorf("help should list editions, got: %s", out)
		}
	})
	t.Run("unknown edition", func(t *testing.T) {
		code, _, errStr := runCapture(t, "mongoose")
		if code != 2 {
			t.Errorf("exit = %d, want 2", code)
		}
		if !strings.Contains(errStr, "unknown edition") {
			t.Errorf("expected unknown-edition message, got: %s", errStr)
		}
	})
}
