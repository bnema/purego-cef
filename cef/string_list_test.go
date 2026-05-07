package cef

import (
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	portoutmocks "github.com/bnema/purego-cef/internal/ports/out/mocks"
)

type fileDialogCallbackRecorder struct {
	contPaths StringList
	canceled  int
}

func (r *fileDialogCallbackRecorder) Cont(filePaths StringList) {
	r.contPaths = filePaths
}

func (r *fileDialogCallbackRecorder) Cancel() {
	r.canceled++
}

func installStringListTestEngine(t *testing.T, m *portoutmocks.MockCAPI) {
	t.Helper()

	prevEng := eng
	prevCurrentRefManager := currentRefManager
	prevRegisteredRefManagers := append([]*core.RefManager(nil), registeredRefManagers...)

	e := core.New(m)
	eng = e
	setCurrentRefManager(e.Refs())

	t.Cleanup(func() {
		eng = prevEng
		currentRefManager = prevCurrentRefManager
		registeredRefManagers = prevRegisteredRefManagers
	})
}

func newStringListTestMock(t *testing.T) *portoutmocks.MockCAPI {
	t.Helper()

	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
	m.EXPECT().StringSet(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(src *uint16, srcLen uintptr, output unsafe.Pointer, _ int32) int32 {
			out := (*capi.CEFStringT)(output)
			n := int(srcLen)
			if n == 0 || src == nil {
				out.Str = nil
				out.Length = 0
				return 1
			}
			buf := make([]uint16, n)
			copy(buf, unsafe.Slice(src, n))
			out.Str = &buf[0]
			out.Length = uintptr(n)
			return 1
		}).Maybe()
	m.EXPECT().StringClear(mock.Anything).Maybe()
	return m
}

func readCEFStringValue(ptr unsafe.Pointer) string {
	cs := (*capi.CEFStringT)(ptr)
	if cs == nil || cs.Str == nil || cs.Length == 0 {
		return ""
	}
	return string(utf16.Decode(unsafe.Slice(cs.Str, cs.Length)))
}

func writeCEFStringValue(ptr unsafe.Pointer, value string) {
	encoded := utf16.Encode([]rune(value))
	out := (*capi.CEFStringT)(ptr)
	if len(encoded) == 0 {
		out.Str = nil
		out.Length = 0
		return
	}
	buf := make([]uint16, len(encoded))
	copy(buf, encoded)
	out.Str = &buf[0]
	out.Length = uintptr(len(buf))
}

func TestNewStringListAndFreeStringList(t *testing.T) {
	m := newStringListTestMock(t)
	installStringListTestEngine(t, m)

	const listHandle = uintptr(0xfeed)
	appended := make([]string, 0, 2)
	m.EXPECT().StringListAlloc().Return(listHandle).Once()
	m.EXPECT().StringListAppend(listHandle, mock.Anything).
		Run(func(_ uintptr, value unsafe.Pointer) {
			appended = append(appended, readCEFStringValue(value))
		}).Twice()
	m.EXPECT().StringListFree(listHandle).Once()

	list := NewStringList("first.txt", "second.txt")
	require.Equal(t, StringList(listHandle), list)
	require.Equal(t, []string{"first.txt", "second.txt"}, appended)

	FreeStringList(list)
}

func TestStringListToSlice(t *testing.T) {
	m := newStringListTestMock(t)
	installStringListTestEngine(t, m)

	const listHandle = uintptr(0xbeef)
	m.EXPECT().StringListSize(listHandle).Return(uintptr(2)).Once()
	m.EXPECT().StringListValue(listHandle, uintptr(0), mock.Anything).
		RunAndReturn(func(_ uintptr, _ uintptr, value unsafe.Pointer) int32 {
			writeCEFStringValue(value, "image/*")
			return 1
		}).Once()
	m.EXPECT().StringListValue(listHandle, uintptr(1), mock.Anything).
		RunAndReturn(func(_ uintptr, _ uintptr, value unsafe.Pointer) int32 {
			writeCEFStringValue(value, ".png;.jpg")
			return 1
		}).Once()

	got := StringListToSlice(StringList(listHandle))
	require.Equal(t, []string{"image/*", ".png;.jpg"}, got)
}

func TestContinueFileDialog(t *testing.T) {
	t.Run("cancel on empty selection", func(t *testing.T) {
		recorder := &fileDialogCallbackRecorder{}
		ContinueFileDialog(recorder)
		require.Equal(t, 1, recorder.canceled)
		require.Zero(t, recorder.contPaths)
	})

	t.Run("continue with selected paths", func(t *testing.T) {
		m := newStringListTestMock(t)
		installStringListTestEngine(t, m)

		const listHandle = uintptr(0x1234)
		m.EXPECT().StringListAlloc().Return(listHandle).Once()
		m.EXPECT().StringListAppend(listHandle, mock.Anything).Once()
		m.EXPECT().StringListFree(listHandle).Once()

		recorder := &fileDialogCallbackRecorder{}
		ContinueFileDialog(recorder, "/tmp/report.txt")

		require.Equal(t, StringList(listHandle), recorder.contPaths)
		require.Zero(t, recorder.canceled)
	})
}
