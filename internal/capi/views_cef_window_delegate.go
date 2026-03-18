package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFWindowDelegateT struct {
	_                            structs.HostLayout
	Base                         CEFPanelDelegateT
	OnWindowCreated              uintptr
	OnWindowClosing              uintptr
	OnWindowDestroyed            uintptr
	OnWindowActivationChanged    uintptr
	OnWindowBoundsChanged        uintptr
	OnWindowFullscreenTransition uintptr
	GetParentWindow              uintptr
	IsWindowModalDialog          uintptr
	GetInitialBounds             uintptr
	GetInitialShowState          uintptr
	IsFrameless                  uintptr
	WithStandardWindowButtons    uintptr
	GetTitlebarHeight            uintptr
	AcceptsFirstMouse            uintptr
	CanResize                    uintptr
	CanMaximize                  uintptr
	CanMinimize                  uintptr
	CanClose                     uintptr
	OnAccelerator                uintptr
	OnKeyEvent                   uintptr
	OnThemeColorsChanged         uintptr
	GetWindowRuntimeStyle        uintptr
	GetLinuxWindowProperties     uintptr
}

func (v *CEFWindowDelegateT) OverrideOnWindowCreated(fn uintptr) { v.OnWindowCreated = fn }

func (v *CEFWindowDelegateT) CallOnWindowCreated(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWindowCreated, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnWindowClosing(fn uintptr) { v.OnWindowClosing = fn }

func (v *CEFWindowDelegateT) CallOnWindowClosing(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWindowClosing, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnWindowDestroyed(fn uintptr) { v.OnWindowDestroyed = fn }

func (v *CEFWindowDelegateT) CallOnWindowDestroyed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWindowDestroyed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnWindowActivationChanged(fn uintptr) {
	v.OnWindowActivationChanged = fn
}

func (v *CEFWindowDelegateT) CallOnWindowActivationChanged(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWindowActivationChanged, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnWindowBoundsChanged(fn uintptr) { v.OnWindowBoundsChanged = fn }

func (v *CEFWindowDelegateT) CallOnWindowBoundsChanged(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWindowBoundsChanged, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnWindowFullscreenTransition(fn uintptr) {
	v.OnWindowFullscreenTransition = fn
}

func (v *CEFWindowDelegateT) CallOnWindowFullscreenTransition(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWindowFullscreenTransition, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideGetParentWindow(fn uintptr) { v.GetParentWindow = fn }

func (v *CEFWindowDelegateT) CallGetParentWindow(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetParentWindow, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideIsWindowModalDialog(fn uintptr) { v.IsWindowModalDialog = fn }

func (v *CEFWindowDelegateT) CallIsWindowModalDialog(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsWindowModalDialog, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideGetInitialBounds(fn uintptr) { v.GetInitialBounds = fn }

func (v *CEFWindowDelegateT) CallGetInitialBounds(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetInitialBounds, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideGetInitialShowState(fn uintptr) { v.GetInitialShowState = fn }

func (v *CEFWindowDelegateT) CallGetInitialShowState(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetInitialShowState, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideIsFrameless(fn uintptr) { v.IsFrameless = fn }

func (v *CEFWindowDelegateT) CallIsFrameless(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsFrameless, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideWithStandardWindowButtons(fn uintptr) {
	v.WithStandardWindowButtons = fn
}

func (v *CEFWindowDelegateT) CallWithStandardWindowButtons(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.WithStandardWindowButtons, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideGetTitlebarHeight(fn uintptr) { v.GetTitlebarHeight = fn }

func (v *CEFWindowDelegateT) CallGetTitlebarHeight(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTitlebarHeight, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideAcceptsFirstMouse(fn uintptr) { v.AcceptsFirstMouse = fn }

func (v *CEFWindowDelegateT) CallAcceptsFirstMouse(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AcceptsFirstMouse, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideCanResize(fn uintptr) { v.CanResize = fn }

func (v *CEFWindowDelegateT) CallCanResize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CanResize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideCanMaximize(fn uintptr) { v.CanMaximize = fn }

func (v *CEFWindowDelegateT) CallCanMaximize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CanMaximize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideCanMinimize(fn uintptr) { v.CanMinimize = fn }

func (v *CEFWindowDelegateT) CallCanMinimize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CanMinimize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideCanClose(fn uintptr) { v.CanClose = fn }

func (v *CEFWindowDelegateT) CallCanClose(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CanClose, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnAccelerator(fn uintptr) { v.OnAccelerator = fn }

func (v *CEFWindowDelegateT) CallOnAccelerator(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAccelerator, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnKeyEvent(fn uintptr) { v.OnKeyEvent = fn }

func (v *CEFWindowDelegateT) CallOnKeyEvent(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnKeyEvent, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideOnThemeColorsChanged(fn uintptr) { v.OnThemeColorsChanged = fn }

func (v *CEFWindowDelegateT) CallOnThemeColorsChanged(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnThemeColorsChanged, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideGetWindowRuntimeStyle(fn uintptr) { v.GetWindowRuntimeStyle = fn }

func (v *CEFWindowDelegateT) CallGetWindowRuntimeStyle(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetWindowRuntimeStyle, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFWindowDelegateT) OverrideGetLinuxWindowProperties(fn uintptr) {
	v.GetLinuxWindowProperties = fn
}

func (v *CEFWindowDelegateT) CallGetLinuxWindowProperties(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLinuxWindowProperties, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterWindowDelegate(handle uintptr) {
}
