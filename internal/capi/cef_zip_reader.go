package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFZipReaderT struct {
	_                   structs.HostLayout
	Base                CEFBaseRefCountedT
	MoveToFirstFile     uintptr
	MoveToNextFile      uintptr
	MoveToFile          uintptr
	Close               uintptr
	GetFileName         uintptr
	GetFileSize         uintptr
	GetFileLastModified uintptr
	OpenFile            uintptr
	CloseFile           uintptr
	ReadFile            uintptr
	Tell                uintptr
	Eof                 uintptr
}

func (v *CEFZipReaderT) OverrideMoveToFirstFile(fn uintptr) { v.MoveToFirstFile = fn }

func (v *CEFZipReaderT) CallMoveToFirstFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.MoveToFirstFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideMoveToNextFile(fn uintptr) { v.MoveToNextFile = fn }

func (v *CEFZipReaderT) CallMoveToNextFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.MoveToNextFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideMoveToFile(fn uintptr) { v.MoveToFile = fn }

func (v *CEFZipReaderT) CallMoveToFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.MoveToFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideClose(fn uintptr) { v.Close = fn }

func (v *CEFZipReaderT) CallClose(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Close, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideGetFileName(fn uintptr) { v.GetFileName = fn }

func (v *CEFZipReaderT) CallGetFileName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFileName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideGetFileSize(fn uintptr) { v.GetFileSize = fn }

func (v *CEFZipReaderT) CallGetFileSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFileSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideGetFileLastModified(fn uintptr) { v.GetFileLastModified = fn }

func (v *CEFZipReaderT) CallGetFileLastModified(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFileLastModified, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideOpenFile(fn uintptr) { v.OpenFile = fn }

func (v *CEFZipReaderT) CallOpenFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OpenFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideCloseFile(fn uintptr) { v.CloseFile = fn }

func (v *CEFZipReaderT) CallCloseFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CloseFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideReadFile(fn uintptr) { v.ReadFile = fn }

func (v *CEFZipReaderT) CallReadFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ReadFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideTell(fn uintptr) { v.Tell = fn }

func (v *CEFZipReaderT) CallTell(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Tell, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFZipReaderT) OverrideEof(fn uintptr) { v.Eof = fn }

func (v *CEFZipReaderT) CallEof(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Eof, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFZipReaderCreate func(Stream unsafe.Pointer) unsafe.Pointer

func RegisterZipReader(handle uintptr) {
	purego.RegisterLibFunc(&CEFZipReaderCreate, handle, "cef_zip_reader_create")
}
