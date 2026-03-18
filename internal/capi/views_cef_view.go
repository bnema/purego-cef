package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFViewT struct {
	_                        structs.HostLayout
	Base                     CEFBaseRefCountedT
	AsBrowserView            uintptr
	AsButton                 uintptr
	AsPanel                  uintptr
	AsScrollView             uintptr
	AsTextfield              uintptr
	GetTypeString            uintptr
	ToString                 uintptr
	IsValid                  uintptr
	IsAttached               uintptr
	IsSame                   uintptr
	GetDelegate              uintptr
	GetWindow                uintptr
	GetID                    uintptr
	SetID                    uintptr
	GetGroupID               uintptr
	SetGroupID               uintptr
	GetParentView            uintptr
	GetViewForID             uintptr
	SetBounds                uintptr
	GetBounds                uintptr
	GetBoundsInScreen        uintptr
	SetSize                  uintptr
	GetSize                  uintptr
	SetPosition              uintptr
	GetPosition              uintptr
	SetInsets                uintptr
	GetInsets                uintptr
	GetPreferredSize         uintptr
	SizeToPreferredSize      uintptr
	GetMinimumSize           uintptr
	GetMaximumSize           uintptr
	GetHeightForWidth        uintptr
	InvalidateLayout         uintptr
	SetVisible               uintptr
	IsVisible                uintptr
	IsDrawn                  uintptr
	SetEnabled               uintptr
	IsEnabled                uintptr
	SetFocusable             uintptr
	IsFocusable              uintptr
	IsAccessibilityFocusable uintptr
	HasFocus                 uintptr
	RequestFocus             uintptr
	SetBackgroundColor       uintptr
	GetBackgroundColor       uintptr
	GetThemeColor            uintptr
	ConvertPointToScreen     uintptr
	ConvertPointFromScreen   uintptr
	ConvertPointToWindow     uintptr
	ConvertPointFromWindow   uintptr
	ConvertPointToView       uintptr
	ConvertPointFromView     uintptr
}

func (v *CEFViewT) OverrideAsBrowserView(fn uintptr) { v.AsBrowserView = fn }

func (v *CEFViewT) CallAsBrowserView(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsBrowserView, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideAsButton(fn uintptr) { v.AsButton = fn }

func (v *CEFViewT) CallAsButton(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsButton, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideAsPanel(fn uintptr) { v.AsPanel = fn }

func (v *CEFViewT) CallAsPanel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsPanel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideAsScrollView(fn uintptr) { v.AsScrollView = fn }

func (v *CEFViewT) CallAsScrollView(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsScrollView, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideAsTextfield(fn uintptr) { v.AsTextfield = fn }

func (v *CEFViewT) CallAsTextfield(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsTextfield, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetTypeString(fn uintptr) { v.GetTypeString = fn }

func (v *CEFViewT) CallGetTypeString(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTypeString, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideToString(fn uintptr) { v.ToString = fn }

func (v *CEFViewT) CallToString(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ToString, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsValid(fn uintptr) { v.IsValid = fn }

func (v *CEFViewT) CallIsValid(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsValid, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsAttached(fn uintptr) { v.IsAttached = fn }

func (v *CEFViewT) CallIsAttached(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsAttached, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsSame(fn uintptr) { v.IsSame = fn }

func (v *CEFViewT) CallIsSame(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsSame, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetDelegate(fn uintptr) { v.GetDelegate = fn }

func (v *CEFViewT) CallGetDelegate(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDelegate, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetWindow(fn uintptr) { v.GetWindow = fn }

func (v *CEFViewT) CallGetWindow(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetWindow, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetID(fn uintptr) { v.GetID = fn }

func (v *CEFViewT) CallGetID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetID(fn uintptr) { v.SetID = fn }

func (v *CEFViewT) CallSetID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetGroupID(fn uintptr) { v.GetGroupID = fn }

func (v *CEFViewT) CallGetGroupID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetGroupID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetGroupID(fn uintptr) { v.SetGroupID = fn }

func (v *CEFViewT) CallSetGroupID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetGroupID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetParentView(fn uintptr) { v.GetParentView = fn }

func (v *CEFViewT) CallGetParentView(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetParentView, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetViewForID(fn uintptr) { v.GetViewForID = fn }

func (v *CEFViewT) CallGetViewForID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetViewForID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetBounds(fn uintptr) { v.SetBounds = fn }

func (v *CEFViewT) CallSetBounds(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetBounds, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetBounds(fn uintptr) { v.GetBounds = fn }

func (v *CEFViewT) CallGetBounds(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetBounds, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetBoundsInScreen(fn uintptr) { v.GetBoundsInScreen = fn }

func (v *CEFViewT) CallGetBoundsInScreen(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetBoundsInScreen, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetSize(fn uintptr) { v.SetSize = fn }

func (v *CEFViewT) CallSetSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetSize(fn uintptr) { v.GetSize = fn }

func (v *CEFViewT) CallGetSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetPosition(fn uintptr) { v.SetPosition = fn }

func (v *CEFViewT) CallSetPosition(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetPosition, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetPosition(fn uintptr) { v.GetPosition = fn }

func (v *CEFViewT) CallGetPosition(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPosition, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetInsets(fn uintptr) { v.SetInsets = fn }

func (v *CEFViewT) CallSetInsets(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetInsets, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetInsets(fn uintptr) { v.GetInsets = fn }

func (v *CEFViewT) CallGetInsets(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetInsets, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetPreferredSize(fn uintptr) { v.GetPreferredSize = fn }

func (v *CEFViewT) CallGetPreferredSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPreferredSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSizeToPreferredSize(fn uintptr) { v.SizeToPreferredSize = fn }

func (v *CEFViewT) CallSizeToPreferredSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SizeToPreferredSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetMinimumSize(fn uintptr) { v.GetMinimumSize = fn }

func (v *CEFViewT) CallGetMinimumSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMinimumSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetMaximumSize(fn uintptr) { v.GetMaximumSize = fn }

func (v *CEFViewT) CallGetMaximumSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMaximumSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetHeightForWidth(fn uintptr) { v.GetHeightForWidth = fn }

func (v *CEFViewT) CallGetHeightForWidth(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHeightForWidth, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideInvalidateLayout(fn uintptr) { v.InvalidateLayout = fn }

func (v *CEFViewT) CallInvalidateLayout(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InvalidateLayout, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetVisible(fn uintptr) { v.SetVisible = fn }

func (v *CEFViewT) CallSetVisible(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetVisible, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsVisible(fn uintptr) { v.IsVisible = fn }

func (v *CEFViewT) CallIsVisible(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsVisible, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsDrawn(fn uintptr) { v.IsDrawn = fn }

func (v *CEFViewT) CallIsDrawn(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsDrawn, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetEnabled(fn uintptr) { v.SetEnabled = fn }

func (v *CEFViewT) CallSetEnabled(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetEnabled, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsEnabled(fn uintptr) { v.IsEnabled = fn }

func (v *CEFViewT) CallIsEnabled(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsEnabled, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetFocusable(fn uintptr) { v.SetFocusable = fn }

func (v *CEFViewT) CallSetFocusable(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetFocusable, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsFocusable(fn uintptr) { v.IsFocusable = fn }

func (v *CEFViewT) CallIsFocusable(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsFocusable, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideIsAccessibilityFocusable(fn uintptr) { v.IsAccessibilityFocusable = fn }

func (v *CEFViewT) CallIsAccessibilityFocusable(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsAccessibilityFocusable, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideHasFocus(fn uintptr) { v.HasFocus = fn }

func (v *CEFViewT) CallHasFocus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasFocus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideRequestFocus(fn uintptr) { v.RequestFocus = fn }

func (v *CEFViewT) CallRequestFocus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RequestFocus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideSetBackgroundColor(fn uintptr) { v.SetBackgroundColor = fn }

func (v *CEFViewT) CallSetBackgroundColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetBackgroundColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetBackgroundColor(fn uintptr) { v.GetBackgroundColor = fn }

func (v *CEFViewT) CallGetBackgroundColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetBackgroundColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideGetThemeColor(fn uintptr) { v.GetThemeColor = fn }

func (v *CEFViewT) CallGetThemeColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetThemeColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideConvertPointToScreen(fn uintptr) { v.ConvertPointToScreen = fn }

func (v *CEFViewT) CallConvertPointToScreen(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ConvertPointToScreen, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideConvertPointFromScreen(fn uintptr) { v.ConvertPointFromScreen = fn }

func (v *CEFViewT) CallConvertPointFromScreen(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ConvertPointFromScreen, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideConvertPointToWindow(fn uintptr) { v.ConvertPointToWindow = fn }

func (v *CEFViewT) CallConvertPointToWindow(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ConvertPointToWindow, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideConvertPointFromWindow(fn uintptr) { v.ConvertPointFromWindow = fn }

func (v *CEFViewT) CallConvertPointFromWindow(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ConvertPointFromWindow, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideConvertPointToView(fn uintptr) { v.ConvertPointToView = fn }

func (v *CEFViewT) CallConvertPointToView(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ConvertPointToView, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFViewT) OverrideConvertPointFromView(fn uintptr) { v.ConvertPointFromView = fn }

func (v *CEFViewT) CallConvertPointFromView(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ConvertPointFromView, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterView(handle uintptr) {
}
