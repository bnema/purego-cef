package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFResolveCallbackT struct {
	_                  structs.HostLayout
	Base               CEFBaseRefCountedT
	OnResolveCompleted uintptr
}

func (v *CEFResolveCallbackT) OverrideOnResolveCompleted(fn uintptr) { v.OnResolveCompleted = fn }

func (v *CEFResolveCallbackT) CallOnResolveCompleted(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnResolveCompleted, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFSettingObserverT struct {
	_                structs.HostLayout
	Base             CEFBaseRefCountedT
	OnSettingChanged uintptr
}

func (v *CEFSettingObserverT) OverrideOnSettingChanged(fn uintptr) { v.OnSettingChanged = fn }

func (v *CEFSettingObserverT) CallOnSettingChanged(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnSettingChanged, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFRequestContextT struct {
	_                            structs.HostLayout
	Base                         CEFPreferenceManagerT
	IsSame                       uintptr
	IsSharingWith                uintptr
	IsGlobal                     uintptr
	GetHandler                   uintptr
	GetCachePath                 uintptr
	GetCookieManager             uintptr
	RegisterSchemeHandlerFactory uintptr
	ClearSchemeHandlerFactories  uintptr
	ClearCertificateExceptions   uintptr
	ClearHttpAuthCredentials     uintptr
	CloseAllConnections          uintptr
	ResolveHost                  uintptr
	GetMediaRouter               uintptr
	GetWebsiteSetting            uintptr
	SetWebsiteSetting            uintptr
	GetContentSetting            uintptr
	SetContentSetting            uintptr
	SetChromeColorScheme         uintptr
	GetChromeColorSchemeMode     uintptr
	GetChromeColorSchemeColor    uintptr
	GetChromeColorSchemeVariant  uintptr
	AddSettingObserver           uintptr
	ClearHttpCache               uintptr
}

func (v *CEFRequestContextT) OverrideIsSame(fn uintptr) { v.IsSame = fn }

func (v *CEFRequestContextT) CallIsSame(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsSame, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideIsSharingWith(fn uintptr) { v.IsSharingWith = fn }

func (v *CEFRequestContextT) CallIsSharingWith(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsSharingWith, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideIsGlobal(fn uintptr) { v.IsGlobal = fn }

func (v *CEFRequestContextT) CallIsGlobal(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsGlobal, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetHandler(fn uintptr) { v.GetHandler = fn }

func (v *CEFRequestContextT) CallGetHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetCachePath(fn uintptr) { v.GetCachePath = fn }

func (v *CEFRequestContextT) CallGetCachePath(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCachePath, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetCookieManager(fn uintptr) { v.GetCookieManager = fn }

func (v *CEFRequestContextT) CallGetCookieManager(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCookieManager, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideRegisterSchemeHandlerFactory(fn uintptr) {
	v.RegisterSchemeHandlerFactory = fn
}

func (v *CEFRequestContextT) CallRegisterSchemeHandlerFactory(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RegisterSchemeHandlerFactory, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideClearSchemeHandlerFactories(fn uintptr) {
	v.ClearSchemeHandlerFactories = fn
}

func (v *CEFRequestContextT) CallClearSchemeHandlerFactories(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ClearSchemeHandlerFactories, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideClearCertificateExceptions(fn uintptr) {
	v.ClearCertificateExceptions = fn
}

func (v *CEFRequestContextT) CallClearCertificateExceptions(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ClearCertificateExceptions, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideClearHttpAuthCredentials(fn uintptr) {
	v.ClearHttpAuthCredentials = fn
}

func (v *CEFRequestContextT) CallClearHttpAuthCredentials(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ClearHttpAuthCredentials, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideCloseAllConnections(fn uintptr) { v.CloseAllConnections = fn }

func (v *CEFRequestContextT) CallCloseAllConnections(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CloseAllConnections, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideResolveHost(fn uintptr) { v.ResolveHost = fn }

func (v *CEFRequestContextT) CallResolveHost(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ResolveHost, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetMediaRouter(fn uintptr) { v.GetMediaRouter = fn }

func (v *CEFRequestContextT) CallGetMediaRouter(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMediaRouter, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetWebsiteSetting(fn uintptr) { v.GetWebsiteSetting = fn }

func (v *CEFRequestContextT) CallGetWebsiteSetting(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetWebsiteSetting, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideSetWebsiteSetting(fn uintptr) { v.SetWebsiteSetting = fn }

func (v *CEFRequestContextT) CallSetWebsiteSetting(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetWebsiteSetting, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetContentSetting(fn uintptr) { v.GetContentSetting = fn }

func (v *CEFRequestContextT) CallGetContentSetting(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetContentSetting, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideSetContentSetting(fn uintptr) { v.SetContentSetting = fn }

func (v *CEFRequestContextT) CallSetContentSetting(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetContentSetting, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideSetChromeColorScheme(fn uintptr) { v.SetChromeColorScheme = fn }

func (v *CEFRequestContextT) CallSetChromeColorScheme(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetChromeColorScheme, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetChromeColorSchemeMode(fn uintptr) {
	v.GetChromeColorSchemeMode = fn
}

func (v *CEFRequestContextT) CallGetChromeColorSchemeMode(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetChromeColorSchemeMode, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetChromeColorSchemeColor(fn uintptr) {
	v.GetChromeColorSchemeColor = fn
}

func (v *CEFRequestContextT) CallGetChromeColorSchemeColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetChromeColorSchemeColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideGetChromeColorSchemeVariant(fn uintptr) {
	v.GetChromeColorSchemeVariant = fn
}

func (v *CEFRequestContextT) CallGetChromeColorSchemeVariant(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetChromeColorSchemeVariant, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideAddSettingObserver(fn uintptr) { v.AddSettingObserver = fn }

func (v *CEFRequestContextT) CallAddSettingObserver(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddSettingObserver, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextT) OverrideClearHttpCache(fn uintptr) { v.ClearHttpCache = fn }

func (v *CEFRequestContextT) CallClearHttpCache(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ClearHttpCache, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFRequestContextGetGlobalContext func() unsafe.Pointer

var CEFRequestContextCreateContext func(Settings unsafe.Pointer, Handler unsafe.Pointer) unsafe.Pointer

var CEFRequestContextCEFCreateContextShared func(Other unsafe.Pointer, Handler unsafe.Pointer) unsafe.Pointer

func RegisterRequestContext(handle uintptr) {
	purego.RegisterLibFunc(&CEFRequestContextGetGlobalContext, handle, "cef_request_context_get_global_context")
	purego.RegisterLibFunc(&CEFRequestContextCreateContext, handle, "cef_request_context_create_context")
	purego.RegisterLibFunc(&CEFRequestContextCEFCreateContextShared, handle, "cef_request_context_cef_create_context_shared")
}
