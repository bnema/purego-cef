package cef

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"

	"github.com/bnema/purego-cef/internal/core"
	portout "github.com/bnema/purego-cef/internal/ports/out"
	portoutmocks "github.com/bnema/purego-cef/internal/ports/out/mocks"
)

func restoreInitTestDeps() func() {
	prevOpen := openCEFLibrary
	prevNewCAPI := newCAPI
	prevProcessArgs := processArgs
	prevExit := exitProcess
	prevStderr := stderrWriter
	return func() {
		openCEFLibrary = prevOpen
		newCAPI = prevNewCAPI
		processArgs = prevProcessArgs
		exitProcess = prevExit
		stderrWriter = prevStderr
	}
}

func TestDefaultSettings(t *testing.T) {
	got := DefaultSettings()
	want := Settings{
		ExternalMessagePump:        true,
		WindowlessRenderingEnabled: true,
		NoSandbox:                  true,
	}
	if got != want {
		t.Fatalf("DefaultSettings() = %#v, want %#v", got, want)
	}
}

func TestSettingsCoreSettingsCopiesFields(t *testing.T) {
	s := Settings{
		CEFDir:                     "/cef",
		LogSeverity:                3,
		MultiThreadedMessageLoop:   true,
		WindowlessRenderingEnabled: true,
		ExternalMessagePump:        true,
		NoSandbox:                  true,
		BrowserSubprocessPath:      "/cef/helper",
		LogFile:                    "/tmp/cef.log",
		InitTraceFile:              "/tmp/trace.log",
		CachePath:                  "/tmp/cache",
		RootCachePath:              "/tmp/root-cache",
	}

	got := s.coreSettings()
	want := core.Settings{
		CEFDir:                     "/cef",
		LogSeverity:                3,
		MultiThreadedMessageLoop:   true,
		WindowlessRenderingEnabled: true,
		ExternalMessagePump:        true,
		NoSandbox:                  true,
		BrowserSubprocessPath:      "/cef/helper",
		LogFile:                    "/tmp/cef.log",
		InitTraceFile:              "/tmp/trace.log",
		CachePath:                  "/tmp/cache",
		RootCachePath:              "/tmp/root-cache",
	}
	if got != want {
		t.Fatalf("Settings.coreSettings() = %#v, want %#v", got, want)
	}
}

func TestExecuteSubprocessReturnsExecutedExitCode(t *testing.T) {
	defer restoreInitTestDeps()()

	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().ExecuteProcess(mock.Anything, unsafe.Pointer(nil), unsafe.Pointer(nil)).Return(int32(7)).Once()

	openCEFLibrary = func(string) (uintptr, error) { return 1, nil }
	newCAPI = func(uintptr) portout.CAPI { return m }
	processArgs = func() []string { return []string{"purego-cef-test"} }

	executed, exitCode, err := ExecuteSubprocess()
	if err != nil {
		t.Fatalf("ExecuteSubprocess() error = %v, want nil", err)
	}
	if !executed {
		t.Fatal("ExecuteSubprocess() executed = false, want true")
	}
	if exitCode != 7 {
		t.Fatalf("ExecuteSubprocess() exitCode = %d, want 7", exitCode)
	}
}

func TestExecuteSubprocessReturnsContinueStatus(t *testing.T) {
	defer restoreInitTestDeps()()

	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().ExecuteProcess(mock.Anything, unsafe.Pointer(nil), unsafe.Pointer(nil)).Return(int32(-1)).Once()

	openCEFLibrary = func(string) (uintptr, error) { return 1, nil }
	newCAPI = func(uintptr) portout.CAPI { return m }
	processArgs = func() []string { return []string{"purego-cef-test"} }

	executed, exitCode, err := ExecuteSubprocess()
	if err != nil {
		t.Fatalf("ExecuteSubprocess() error = %v, want nil", err)
	}
	if executed {
		t.Fatal("ExecuteSubprocess() executed = true, want false")
	}
	if exitCode != 0 {
		t.Fatalf("ExecuteSubprocess() exitCode = %d, want 0", exitCode)
	}
}

func TestExecuteSubprocessPropagatesLoaderError(t *testing.T) {
	defer restoreInitTestDeps()()

	wantErr := errors.New("open failed")
	openCEFLibrary = func(string) (uintptr, error) { return 0, wantErr }
	newCAPI = func(uintptr) portout.CAPI {
		t.Fatal("newCAPI should not be called when openCEFLibrary fails")
		return nil
	}

	executed, exitCode, err := ExecuteSubprocess()
	if !errors.Is(err, wantErr) {
		t.Fatalf("ExecuteSubprocess() error = %v, want %v", err, wantErr)
	}
	if executed {
		t.Fatal("ExecuteSubprocess() executed = true, want false")
	}
	if exitCode != 0 {
		t.Fatalf("ExecuteSubprocess() exitCode = %d, want 0", exitCode)
	}
}

func TestMaybeExitSubprocessExitsOnHandledSubprocess(t *testing.T) {
	defer restoreInitTestDeps()()

	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().ExecuteProcess(mock.Anything, unsafe.Pointer(nil), unsafe.Pointer(nil)).Return(int32(11)).Once()

	openCEFLibrary = func(string) (uintptr, error) { return 1, nil }
	newCAPI = func(uintptr) portout.CAPI { return m }
	processArgs = func() []string { return []string{"purego-cef-test"} }

	exited := false
	exitCode := -1
	exitProcess = func(code int) {
		exited = true
		exitCode = code
	}

	MaybeExitSubprocess()
	if !exited {
		t.Fatal("MaybeExitSubprocess() did not exit")
	}
	if exitCode != 11 {
		t.Fatalf("MaybeExitSubprocess() exitCode = %d, want 11", exitCode)
	}
}

func TestMaybeExitSubprocessReportsLoaderError(t *testing.T) {
	defer restoreInitTestDeps()()

	openCEFLibrary = func(string) (uintptr, error) { return 0, errors.New("boom") }
	var stderr bytes.Buffer
	stderrWriter = &stderr
	exitProcess = func(int) {
		t.Fatal("exitProcess should not be called when ExecuteSubprocess fails")
	}

	MaybeExitSubprocess()
	if got := stderr.String(); !strings.Contains(got, "ExecuteSubprocess: boom") {
		t.Fatalf("stderr = %q, want ExecuteSubprocess error message", got)
	}
}
