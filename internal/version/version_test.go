package version

import (
	"os"
	"strings"
	"testing"
)

func capture(f func()) string {
	originalStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = originalStdout

	var buf = make([]byte, 5096)
	n, _ := r.Read(buf)
	return string(buf[:n])
}

func TestRelease(t *testing.T) {
	expected := "1.0.0"
	Version = expected

	if expected != Release() {
		t.Errorf("version should be %s, got %s", expected, Release())
	}

	Version = ""

	if Release() != "dev" {
		t.Errorf("version should be 'dev' when empty, got %s", Release())
	}
}

func TestCommit(t *testing.T) {
	expected := "abc123def"
	GitCommit = expected

	if expected != Commit() {
		t.Errorf("commit hash should be %s, got %s", expected, Commit())
	}
}

func TestBanner(t *testing.T) {
	if Banner() == "" {
		t.Error("banner should not be empty")
	}

	if len(Banner()) != 212 {
		t.Errorf("banner length should be 212, got %d", len(Banner()))
	}
}

func TestPrint(t *testing.T) {
	Version = "1.0.0"
	GitCommit = "abc123def"

	output := capture(Print)

	if output == "" {
		t.Error("output should not be empty")
	}
	if !strings.Contains(output, Banner()) {
		t.Error("output should contain banner")
	}
	if !strings.Contains(output, "Release: 1.0.0") {
		t.Error("output should contain release info")
	}
	if !strings.Contains(output, "Commit:  abc123def") {
		t.Error("output should contain commit info")
	}
}
