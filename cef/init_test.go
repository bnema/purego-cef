package cef

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	portout "github.com/bnema/purego-cef/internal/ports/out"
	portoutmocks "github.com/bnema/purego-cef/internal/ports/out/mocks"
)

type subprocessAppStub struct{}

func (subprocessAppStub) OnBeforeCommandLineProcessing(string, CommandLine) {}
func (subprocessAppStub) OnRegisterCustomSchemes(SchemeRegistrar)           {}
func (subprocessAppStub) GetResourceBundleHandler() ResourceBundleHandler   { return nil }
func (subprocessAppStub) GetBrowserProcessHandler() BrowserProcessHandler   { return nil }
func (subprocessAppStub) GetRenderProcessHandler() RenderProcessHandler     { return nil }

type subprocessCallbackAppStub struct {
	commandLineCalls int
	schemeCalls      int
}

func (a *subprocessCallbackAppStub) OnBeforeCommandLineProcessing(_ string, commandLine CommandLine) {
	a.commandLineCalls++
	commandLine.AppendSwitch("from-subprocess-test")
}

func (a *subprocessCallbackAppStub) OnRegisterCustomSchemes(registrar SchemeRegistrar) {
	a.schemeCalls++
	registrar.AddCustomScheme("subprocess-test", 0)
}

func (*subprocessCallbackAppStub) GetResourceBundleHandler() ResourceBundleHandler { return nil }
func (*subprocessCallbackAppStub) GetBrowserProcessHandler() BrowserProcessHandler { return nil }
func (*subprocessCallbackAppStub) GetRenderProcessHandler() RenderProcessHandler   { return nil }

func restoreInitTestDeps() func() {
	prevOpen := openCEFLibrary
	prevNewCAPI := newCAPI
	prevProcessArgs := processArgs
	prevExit := exitProcess
	prevStderr := stderrWriter
	prevEng := eng
	prevCurrentRefManager := currentRefManager
	prevRegisteredRefManagers := append([]*core.RefManager(nil), registeredRefManagers...)
	return func() {
		openCEFLibrary = prevOpen
		newCAPI = prevNewCAPI
		processArgs = prevProcessArgs
		exitProcess = prevExit
		stderrWriter = prevStderr
		eng = prevEng
		currentRefManager = prevCurrentRefManager
		registeredRefManagers = prevRegisteredRefManagers
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
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
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
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
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

func TestExecuteSubprocessWithAppPassesWrappedAppAndLeavesEngineUntouched(t *testing.T) {
	defer restoreInitTestDeps()()

	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
	m.EXPECT().ExecuteProcess(mock.Anything, mock.Anything, unsafe.Pointer(nil)).
		RunAndReturn(func(_ unsafe.Pointer, application unsafe.Pointer, _ unsafe.Pointer) int32 {
			if application == nil {
				t.Fatal("ExecuteProcess application = nil, want wrapped App pointer")
			}
			return 9
		}).Once()

	openCEFLibrary = func(string) (uintptr, error) { return 1, nil }
	newCAPI = func(uintptr) portout.CAPI { return m }
	processArgs = func() []string { return []string{"purego-cef-test"} }

	sentinelEng := new(core.Engine)
	eng = sentinelEng
	initialManagers := len(registeredRefManagers)

	executed, exitCode, err := ExecuteSubprocessWithApp(subprocessAppStub{})
	if err != nil {
		t.Fatalf("ExecuteSubprocessWithApp(...) error = %v, want nil", err)
	}
	if !executed {
		t.Fatal("ExecuteSubprocessWithApp(...) executed = false, want true")
	}
	if exitCode != 9 {
		t.Fatalf("ExecuteSubprocessWithApp(...) exitCode = %d, want 9", exitCode)
	}
	if eng != sentinelEng {
		t.Fatal("ExecuteSubprocessWithApp(...) mutated global eng")
	}
	if currentRefManager != nil {
		t.Fatal("ExecuteSubprocessWithApp(...) left a temporary current ref manager behind")
	}
	if got := len(registeredRefManagers); got != initialManagers {
		t.Fatalf("registeredRefManagers length = %d, want %d after cleanup", got, initialManagers)
	}
}

func TestExecuteSubprocessWithAppKeepsEngineAvailableForAppCallbacks(t *testing.T) {
	defer restoreInitTestDeps()()

	app := &subprocessCallbackAppStub{}
	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
	m.EXPECT().StringSet(mock.Anything, mock.Anything, mock.Anything, int32(1)).Return(int32(1)).Maybe()
	m.EXPECT().StringClear(mock.Anything).Return().Maybe()
	m.EXPECT().ExecuteProcess(mock.Anything, mock.Anything, unsafe.Pointer(nil)).
		RunAndReturn(func(_ unsafe.Pointer, application unsafe.Pointer, _ unsafe.Pointer) int32 {
			if application == nil {
				t.Fatal("ExecuteProcess application = nil, want wrapped App pointer")
			}
			rawApp := (*capi.CEFAppT)(application)
			rawCommandLine := &capi.CEFCommandLineT{}
			rawRegistrar := &capi.CEFSchemeRegistrarT{}
			rawApp.CallOnBeforeCommandLineProcessing(0, uintptr(unsafe.Pointer(rawCommandLine)))
			rawApp.CallOnRegisterCustomSchemes(uintptr(unsafe.Pointer(rawRegistrar)))
			return 5
		}).Once()

	openCEFLibrary = func(string) (uintptr, error) { return 1, nil }
	newCAPI = func(uintptr) portout.CAPI { return m }
	processArgs = func() []string { return []string{"purego-cef-test"} }

	executed, exitCode, err := ExecuteSubprocessWithApp(app)
	if err != nil {
		t.Fatalf("ExecuteSubprocessWithApp(...) error = %v, want nil", err)
	}
	if !executed {
		t.Fatal("ExecuteSubprocessWithApp(...) executed = false, want true")
	}
	if exitCode != 5 {
		t.Fatalf("ExecuteSubprocessWithApp(...) exitCode = %d, want 5", exitCode)
	}
	if app.commandLineCalls != 1 {
		t.Fatalf("OnBeforeCommandLineProcessing call count = %d, want 1", app.commandLineCalls)
	}
	if app.schemeCalls != 1 {
		t.Fatalf("OnRegisterCustomSchemes call count = %d, want 1", app.schemeCalls)
	}
}

func TestMaybeExitSubprocessExitsOnHandledSubprocess(t *testing.T) {
	defer restoreInitTestDeps()()

	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
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
