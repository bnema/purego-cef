package capi

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

var CEFResolveURL func(BaseURL unsafe.Pointer, RelativeURL unsafe.Pointer, ResolvedURL unsafe.Pointer) int32

var CEFParseURL func(URL unsafe.Pointer, Parts unsafe.Pointer) int32

var CEFCreateURL func(Parts unsafe.Pointer, URL unsafe.Pointer) int32

var CEFFormatURLForSecurityDisplay func(OriginURL unsafe.Pointer) uintptr

var CEFGetMimeType func(Extension unsafe.Pointer) uintptr

var CEFGetExtensionsForMimeType func(MimeType unsafe.Pointer, Extensions uintptr)

var CEFBase64Encode func(Data unsafe.Pointer, DataSize uintptr) uintptr

var CEFBase64Decode func(Data unsafe.Pointer) unsafe.Pointer

var CEFUriencode func(Text unsafe.Pointer, UsePlus int32) uintptr

var CEFUridecode func(Text unsafe.Pointer, ConvertToUtf8 int32, UnescapeRule CEFUriUnescapeRuleT) uintptr

var CEFParseJson func(JsonString unsafe.Pointer, Options CEFJsonParserOptionsT) unsafe.Pointer

var CEFParseJsonBuffer func(Json unsafe.Pointer, JsonSize uintptr, Options CEFJsonParserOptionsT) unsafe.Pointer

var CEFParseJsonandReturnError func(JsonString unsafe.Pointer, Options CEFJsonParserOptionsT, ErrorMsgOut unsafe.Pointer) unsafe.Pointer

var CEFWriteJson func(Node unsafe.Pointer, Options CEFJsonWriterOptionsT) uintptr

func RegisterParser(handle uintptr) {
	purego.RegisterLibFunc(&CEFResolveURL, handle, "cef_resolve_url")
	purego.RegisterLibFunc(&CEFParseURL, handle, "cef_parse_url")
	purego.RegisterLibFunc(&CEFCreateURL, handle, "cef_create_url")
	purego.RegisterLibFunc(&CEFFormatURLForSecurityDisplay, handle, "cef_format_url_for_security_display")
	purego.RegisterLibFunc(&CEFGetMimeType, handle, "cef_get_mime_type")
	purego.RegisterLibFunc(&CEFGetExtensionsForMimeType, handle, "cef_get_extensions_for_mime_type")
	purego.RegisterLibFunc(&CEFBase64Encode, handle, "cef_base64_encode")
	purego.RegisterLibFunc(&CEFBase64Decode, handle, "cef_base64_decode")
	purego.RegisterLibFunc(&CEFUriencode, handle, "cef_uriencode")
	purego.RegisterLibFunc(&CEFUridecode, handle, "cef_uridecode")
	purego.RegisterLibFunc(&CEFParseJson, handle, "cef_parse_json")
	purego.RegisterLibFunc(&CEFParseJsonBuffer, handle, "cef_parse_json_buffer")
	purego.RegisterLibFunc(&CEFParseJsonandReturnError, handle, "cef_parse_jsonand_return_error")
	purego.RegisterLibFunc(&CEFWriteJson, handle, "cef_write_json")
}
