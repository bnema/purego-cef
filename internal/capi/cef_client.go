package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFClientT struct {
	_                        structs.HostLayout
	Base                     uintptr
	GetAudioHandler          uintptr
	GetCommandHandler        uintptr
	GetContextMenuHandler    uintptr
	GetDialogHandler         uintptr
	GetDisplayHandler        uintptr
	GetDownloadHandler       uintptr
	GetDragHandler           uintptr
	GetFindHandler           uintptr
	GetFocusHandler          uintptr
	GetFrameHandler          uintptr
	GetPermissionHandler     uintptr
	GetJsdialogHandler       uintptr
	GetKeyboardHandler       uintptr
	GetLifeSpanHandler       uintptr
	GetLoadHandler           uintptr
	GetPrintHandler          uintptr
	GetRenderHandler         uintptr
	GetRequestHandler        uintptr
	OnProcessMessageReceived uintptr
}

func (v *CEFClientT) OverrideGetAudioHandler(fn uintptr) { v.GetAudioHandler = fn }

func (v *CEFClientT) CallGetAudioHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAudioHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetCommandHandler(fn uintptr) { v.GetCommandHandler = fn }

func (v *CEFClientT) CallGetCommandHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCommandHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetContextMenuHandler(fn uintptr) { v.GetContextMenuHandler = fn }

func (v *CEFClientT) CallGetContextMenuHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetContextMenuHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetDialogHandler(fn uintptr) { v.GetDialogHandler = fn }

func (v *CEFClientT) CallGetDialogHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDialogHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetDisplayHandler(fn uintptr) { v.GetDisplayHandler = fn }

func (v *CEFClientT) CallGetDisplayHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDisplayHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetDownloadHandler(fn uintptr) { v.GetDownloadHandler = fn }

func (v *CEFClientT) CallGetDownloadHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDownloadHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetDragHandler(fn uintptr) { v.GetDragHandler = fn }

func (v *CEFClientT) CallGetDragHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDragHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetFindHandler(fn uintptr) { v.GetFindHandler = fn }

func (v *CEFClientT) CallGetFindHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFindHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetFocusHandler(fn uintptr) { v.GetFocusHandler = fn }

func (v *CEFClientT) CallGetFocusHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFocusHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetFrameHandler(fn uintptr) { v.GetFrameHandler = fn }

func (v *CEFClientT) CallGetFrameHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFrameHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetPermissionHandler(fn uintptr) { v.GetPermissionHandler = fn }

func (v *CEFClientT) CallGetPermissionHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPermissionHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetJsdialogHandler(fn uintptr) { v.GetJsdialogHandler = fn }

func (v *CEFClientT) CallGetJsdialogHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetJsdialogHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetKeyboardHandler(fn uintptr) { v.GetKeyboardHandler = fn }

func (v *CEFClientT) CallGetKeyboardHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetKeyboardHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetLifeSpanHandler(fn uintptr) { v.GetLifeSpanHandler = fn }

func (v *CEFClientT) CallGetLifeSpanHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLifeSpanHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetLoadHandler(fn uintptr) { v.GetLoadHandler = fn }

func (v *CEFClientT) CallGetLoadHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLoadHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetPrintHandler(fn uintptr) { v.GetPrintHandler = fn }

func (v *CEFClientT) CallGetPrintHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPrintHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetRenderHandler(fn uintptr) { v.GetRenderHandler = fn }

func (v *CEFClientT) CallGetRenderHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetRenderHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideGetRequestHandler(fn uintptr) { v.GetRequestHandler = fn }

func (v *CEFClientT) CallGetRequestHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetRequestHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFClientT) OverrideOnProcessMessageReceived(fn uintptr) { v.OnProcessMessageReceived = fn }

func (v *CEFClientT) CallOnProcessMessageReceived(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnProcessMessageReceived, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterClient(handle uintptr) {
}
