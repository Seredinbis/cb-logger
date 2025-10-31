package logger

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func resetLogger() {
	AppLogger = Logger{}
}

var ansiRegexp = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string { return ansiRegexp.ReplaceAllString(s, "") }

func TestInitAndInfoLogging(t *testing.T) {
	defer resetLogger()

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe err: %v", err)
	}
	os.Stdout = w
	defer func() { _ = w.Close(); os.Stdout = oldStdout; _ = r.Close() }()

	_ = AppLogger.Init()

	Infof("hello %s", "world")

	_ = w.Close()
	_, _ = io.Copy(&buf, r)
	out := stripANSI(buf.String())
	if !strings.Contains(out, "hello world") {
		t.Fatalf("expected message in output, got: %q", out)
	}
}

func TestLoggingBeforeInitDoesNotPanic(t *testing.T) {
	defer resetLogger()

	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	Fatal("fatal")
}

func TestSingletonInitOnlyOnce(t *testing.T) {
	defer resetLogger()

	var buf1 bytes.Buffer
	oldStdout := os.Stdout
	r1, w1, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe1 err: %v", err)
	}
	os.Stdout = w1
	defer func() { _ = w1.Close(); os.Stdout = oldStdout; _ = r1.Close() }()
	_ = AppLogger.Init()

	_ = AppLogger.Init()

	Info("message")
	_ = w1.Close()
	_, _ = io.Copy(&buf1, r1)
	if buf1.Len() == 0 {
		t.Fatalf("expected first writer to receive logs")
	}
}

func TestLevelFiltering(t *testing.T) {
	defer resetLogger()

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe err: %v", err)
	}
	os.Stdout = w
	defer func() { _ = w.Close(); os.Stdout = oldStdout; _ = r.Close() }()
	_ = AppLogger.Init()

	Debug("hidden debug")
	Info("shown info")

	_ = w.Close()
	_, _ = io.Copy(&buf, r)
	out := buf.String()
	if strings.Contains(out, "hidden debug") {
		t.Fatalf("did not expect debug message at info level; got: %q", out)
	}
	if !strings.Contains(out, "shown info") {
		t.Fatalf("expected info message in output; got: %q", out)
	}
}

func TestWarnAndErrorLogging(t *testing.T) {
	defer resetLogger()

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe err: %v", err)
	}
	os.Stdout = w
	defer func() { _ = w.Close(); os.Stdout = oldStdout; _ = r.Close() }()
	_ = AppLogger.Init()

	Warn("warn message")
	Error("error message")

	_ = w.Close()
	_, _ = io.Copy(&buf, r)
	out := stripANSI(buf.String())
	if !strings.Contains(out, "warn message") {
		t.Fatalf("expected warn message; got: %q", out)
	}
	if !strings.Contains(out, "error message") {
		t.Fatalf("expected error message; got: %q", out)
	}
}

func TestInitUsesStdoutWhenWriterNil(t *testing.T) {
	defer resetLogger()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = oldStdout
		_ = r.Close()
	}()

	_ = AppLogger.Init()
	Info("to stdout")

	_ = w.Close()
	var captured bytes.Buffer
	_, _ = io.Copy(&captured, r)

	out := stripANSI(captured.String())
	if !strings.Contains(out, "to stdout") {
		t.Fatalf("expected output on stdout; got: %q", out)
	}
}

func TestFatalExitsProcess(t *testing.T) {
	defer resetLogger()

	cmd := exec.Command(os.Args[0], "-test.run", "TestFatalHelper")
	cmd.Env = append(os.Environ(), "CB_FATAL_TEST=1")
	out, err := cmd.CombinedOutput()
	output := stripANSI(string(out))

	if err == nil {
		t.Fatalf("expected subprocess to exit non-zero due to Fatal; output: %q", output)
	}
	if !strings.Contains(output, "fatal subprocess") {
		t.Fatalf("expected fatal message in output; got: %q", output)
	}
}

func TestFatalHelper(t *testing.T) {
	if os.Getenv("CB_FATAL_TEST") != "1" {
		t.Skip("helper subprocess only")
		return
	}
	_ = AppLogger.Init()
	Fatalf("fatal %s", "subprocess")
	os.Exit(0)
}
