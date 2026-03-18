package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFX509CertPrincipalT struct {
	_                        structs.HostLayout
	Base                     CEFBaseRefCountedT
	GetDisplayName           uintptr
	GetCommonName            uintptr
	GetLocalityName          uintptr
	GetStateOrProvinceName   uintptr
	GetCountryName           uintptr
	GetOrganizationNames     uintptr
	GetOrganizationUnitNames uintptr
}

func (v *CEFX509CertPrincipalT) OverrideGetDisplayName(fn uintptr) { v.GetDisplayName = fn }

func (v *CEFX509CertPrincipalT) CallGetDisplayName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDisplayName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertPrincipalT) OverrideGetCommonName(fn uintptr) { v.GetCommonName = fn }

func (v *CEFX509CertPrincipalT) CallGetCommonName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCommonName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertPrincipalT) OverrideGetLocalityName(fn uintptr) { v.GetLocalityName = fn }

func (v *CEFX509CertPrincipalT) CallGetLocalityName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLocalityName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertPrincipalT) OverrideGetStateOrProvinceName(fn uintptr) {
	v.GetStateOrProvinceName = fn
}

func (v *CEFX509CertPrincipalT) CallGetStateOrProvinceName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetStateOrProvinceName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertPrincipalT) OverrideGetCountryName(fn uintptr) { v.GetCountryName = fn }

func (v *CEFX509CertPrincipalT) CallGetCountryName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCountryName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertPrincipalT) OverrideGetOrganizationNames(fn uintptr) { v.GetOrganizationNames = fn }

func (v *CEFX509CertPrincipalT) CallGetOrganizationNames(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetOrganizationNames, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertPrincipalT) OverrideGetOrganizationUnitNames(fn uintptr) {
	v.GetOrganizationUnitNames = fn
}

func (v *CEFX509CertPrincipalT) CallGetOrganizationUnitNames(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetOrganizationUnitNames, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFX509CertificateT struct {
	_                        structs.HostLayout
	Base                     CEFBaseRefCountedT
	GetSubject               uintptr
	GetIssuer                uintptr
	GetSerialNumber          uintptr
	GetValidStart            uintptr
	GetValidExpiry           uintptr
	GetDerencoded            uintptr
	GetPemencoded            uintptr
	GetIssuerChainSize       uintptr
	GetDerencodedIssuerChain uintptr
	GetPemencodedIssuerChain uintptr
}

func (v *CEFX509CertificateT) OverrideGetSubject(fn uintptr) { v.GetSubject = fn }

func (v *CEFX509CertificateT) CallGetSubject(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSubject, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetIssuer(fn uintptr) { v.GetIssuer = fn }

func (v *CEFX509CertificateT) CallGetIssuer(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetIssuer, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetSerialNumber(fn uintptr) { v.GetSerialNumber = fn }

func (v *CEFX509CertificateT) CallGetSerialNumber(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSerialNumber, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetValidStart(fn uintptr) { v.GetValidStart = fn }

func (v *CEFX509CertificateT) CallGetValidStart(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetValidStart, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetValidExpiry(fn uintptr) { v.GetValidExpiry = fn }

func (v *CEFX509CertificateT) CallGetValidExpiry(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetValidExpiry, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetDerencoded(fn uintptr) { v.GetDerencoded = fn }

func (v *CEFX509CertificateT) CallGetDerencoded(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDerencoded, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetPemencoded(fn uintptr) { v.GetPemencoded = fn }

func (v *CEFX509CertificateT) CallGetPemencoded(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPemencoded, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetIssuerChainSize(fn uintptr) { v.GetIssuerChainSize = fn }

func (v *CEFX509CertificateT) CallGetIssuerChainSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetIssuerChainSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetDerencodedIssuerChain(fn uintptr) {
	v.GetDerencodedIssuerChain = fn
}

func (v *CEFX509CertificateT) CallGetDerencodedIssuerChain(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDerencodedIssuerChain, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFX509CertificateT) OverrideGetPemencodedIssuerChain(fn uintptr) {
	v.GetPemencodedIssuerChain = fn
}

func (v *CEFX509CertificateT) CallGetPemencodedIssuerChain(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPemencodedIssuerChain, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterX509Certificate(handle uintptr) {
}
