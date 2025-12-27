//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type vectorsV1 struct {
	Version   int      `json:"version"`
	UFVK      string   `json:"ufvk"`
	Addresses []string `json:"addresses"`
}

func loadVectors(t *testing.T) vectorsV1 {
	t.Helper()

	root := filepath.Join("..", "..")
	b, err := os.ReadFile(filepath.Join(root, "vectors", "v1.json"))
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}
	var v vectorsV1
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}
	if v.Version != 1 || v.UFVK == "" || len(v.Addresses) != 100 {
		t.Fatalf("invalid vectors")
	}
	return v
}

func run(t *testing.T, bin string, args ...string) (stdout string, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(bin, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err == nil {
		return outBuf.String(), errBuf.String(), 0
	}
	var ee *exec.ExitError
	if !os.IsNotExist(err) && execErrorAs(err, &ee) {
		return outBuf.String(), errBuf.String(), ee.ExitCode()
	}
	t.Fatalf("run %s %s: %v", bin, strings.Join(args, " "), err)
	return "", "", 0
}

func execErrorAs(err error, target **exec.ExitError) bool {
	e, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	*target = e
	return true
}

func TestCLI_DeriveAndBatch(t *testing.T) {
	v := loadVectors(t)

	bin := filepath.Join("..", "..", "bin", "juno-addrgen")
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("missing binary: %v", err)
	}

	stdout, stderr, code := run(t, bin, "derive", "--ufvk", v.UFVK, "--index", "0")
	if code != 0 || stderr != "" {
		t.Fatalf("derive failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
	}
	if got := strings.TrimSpace(stdout); got != v.Addresses[0] {
		t.Fatalf("derive mismatch\nwant: %s\ngot:  %s", v.Addresses[0], got)
	}

	stdout, stderr, code = run(t, bin, "derive", "--ufvk", v.UFVK, "--index", "0", "--json")
	if code != 0 || stderr != "" {
		t.Fatalf("derive json failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
	}
	var resp map[string]any
	if err := json.Unmarshal([]byte(stdout), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["status"] != "ok" || resp["address"] != v.Addresses[0] {
		t.Fatalf("unexpected json: %v", resp)
	}

	stdout, stderr, code = run(t, bin, "batch", "--ufvk", v.UFVK, "--start", "0", "--count", "2")
	if code != 0 || stderr != "" {
		t.Fatalf("batch failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 2 || lines[0] != v.Addresses[0] || lines[1] != v.Addresses[1] {
		t.Fatalf("batch mismatch: %q", stdout)
	}
}
