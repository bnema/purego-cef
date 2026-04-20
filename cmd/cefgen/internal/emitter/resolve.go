package emitter

import (
	"strings"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

// TypeRegistry maps C type names to their parsed metadata and provides
// type resolution from C types to public Go API types.
type TypeRegistry struct {
	structs map[string]*model.Struct
	enums   map[string]*model.Enum
}

// NewTypeRegistry builds a TypeRegistry from all parsed headers.
// Struct and enum types are indexed by their C name for fast lookup.
func NewTypeRegistry(headers []*model.Header) *TypeRegistry {
	r := &TypeRegistry{
		structs: make(map[string]*model.Struct),
		enums:   make(map[string]*model.Enum),
	}
	for _, h := range headers {
		for i := range h.Structs {
			s := &h.Structs[i]
			r.structs[s.CName] = s
		}
		for i := range h.Enums {
			e := &h.Enums[i]
			r.enums[e.CName] = e
		}
	}
	return r
}

// normalizeConst rewrites C type qualifiers so that "const" always appears as a
// leading prefix: "cef_rect_t const*" → "const cef_rect_t*".  This ensures the
// rest of the resolver pipeline can rely on a single canonical form.
func normalizeConst(ct string) string {
	// Strip trailing pointer(s), check for trailing "const", then reassemble.
	starCount := 0
	s := ct
	for {
		trimmed, ok := strings.CutSuffix(s, "*")
		if !ok {
			break
		}
		starCount++
		s = strings.TrimSpace(trimmed)
	}
	if base, ok := strings.CutSuffix(s, " const"); ok {
		return "const " + base + strings.Repeat("*", starCount)
	}
	return ct
}

// ResolvePublicType converts a C type string to the corresponding Go public API type.
func (r *TypeRegistry) ResolvePublicType(ctype string) string {
	ct := normalizeConst(strings.TrimSpace(ctype))

	// Exact matches first.
	switch ct {
	case "void":
		return ""
	case "void*", "const void*":
		return "unsafe.Pointer"
	case "const cef_string_t*":
		return "string"
	case "cef_string_t*":
		return "uintptr"
	case "cef_string_userfree_t":
		return "string"
	case "cef_string_list_t":
		return "StringList"
	case "cef_string_map_t":
		return "StringMap"
	case "cef_string_multimap_t":
		return "StringMultimap"
	case "int":
		return "int32"
	case "unsigned int":
		return "uint32"
	case "size_t":
		return "int"
	case "double":
		return "float64"
	case "float":
		return "float32"
	case "int8_t":
		return "int8"
	case "uint8_t":
		return "uint8"
	case "int16_t":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "int32_t":
		return "int32"
	case "uint32_t":
		return "uint32"
	case "int64_t":
		return "int64"
	case "uint64_t":
		return "uint64"
	case "const char*":
		return "string"
	case "char*":
		return "string"
	case "size_t*":
		return "*int"
	case "int*":
		return "unsafe.Pointer"
	case "uint32_t*":
		return "unsafe.Pointer"
	case "int64_t*":
		return "unsafe.Pointer"
	case "uint64_t*":
		return "unsafe.Pointer"
	case "double*":
		return "unsafe.Pointer"
	case "float*":
		return "unsafe.Pointer"
	}

	// Enum lookup: bare enum type name like "cef_process_id_t".
	if _, ok := r.enums[ct]; ok {
		return publicTypeNameForCName(ct)
	}

	// Handle "struct _cef_xxx_t*" and "const struct _cef_xxx_t*".
	if s, ok := r.resolveStructPointer(ct); ok {
		return s
	}

	// Handle bare "const cef_xxx_t*" pointers (no struct keyword).
	// CEF types like cef_key_event_t are typedefs for struct _cef_key_event_t.
	if s, ok := r.resolveBarePointer(ct); ok {
		return s
	}

	// Enum pointer (rare but possible).
	trimmed := strings.TrimSuffix(ct, "*")
	trimmed = strings.TrimSpace(trimmed)
	if _, ok := r.enums[trimmed]; ok {
		return "*" + publicTypeNameForCName(trimmed)
	}

	return "uintptr"
}

// resolveStructPointer handles "struct _cef_xxx_t*", "const struct _cef_xxx_t*",
// and double-pointer "struct _cef_xxx_t**" forms.
// It strips const, struct prefix, and pointer suffix, then looks up in the registry.
func (r *TypeRegistry) resolveStructPointer(ct string) (string, bool) {
	s := ct
	s = strings.TrimPrefix(s, "const ")
	if !strings.HasPrefix(s, "struct _") {
		return "", false
	}

	// Count and strip pointer stars.
	ptrCount := strings.Count(s, "*")
	s = strings.TrimRight(s, "*")
	s = strings.TrimSpace(s)
	// "struct _cef_browser_t" -> "_cef_browser_t"
	s = strings.TrimPrefix(s, "struct ")

	// For double pointers (output params like **client), use unsafe.Pointer.
	if ptrCount >= 2 {
		return "unsafe.Pointer", true
	}

	// Look up in registry using the underscore-prefixed name.
	if st, ok := r.structs[s]; ok {
		switch st.Kind {
		case "handler", "object":
			return publicTypeNameForCName(st.CName), true
		case "data":
			return "*" + publicTypeNameForCName(s), true
		}
	}

	// Also try without leading underscore: _cef_xxx_t -> cef_xxx_t.
	trimmed := strings.TrimLeft(s, "_")
	if st, ok := r.structs[trimmed]; ok {
		switch st.Kind {
		case "handler", "object":
			return publicTypeNameForCName(st.CName), true
		case "data":
			return "*" + publicTypeNameForCName(trimmed), true
		}
	}

	// Not found in registry; fall back to public name.
	return publicTypeNameForCName(s), true
}

// resolveBarePointer handles "const cef_xxx_t*" or "cef_xxx_t*" patterns
// where the type is a typedef (not written with struct keyword). It strips
// const and pointer, then looks up the bare name with underscore prefix
// in the struct registry.
func (r *TypeRegistry) resolveBarePointer(ct string) (string, bool) {
	s := ct
	s = strings.TrimPrefix(s, "const ")
	if !strings.HasSuffix(s, "*") {
		return "", false
	}
	ptrCount := strings.Count(s, "*")
	s = strings.TrimRight(s, "*")
	s = strings.TrimSpace(s)

	if ptrCount >= 2 {
		return "unsafe.Pointer", true
	}

	// Try with underscore prefix: cef_key_event_t → _cef_key_event_t
	withUnderscore := "_" + s
	if st, ok := r.structs[withUnderscore]; ok {
		switch st.Kind {
		case "handler", "object":
			return publicTypeNameForCName(st.CName), true
		case "data":
			return "*" + publicTypeNameForCName(withUnderscore), true
		}
	}
	// Try bare name as-is.
	if st, ok := r.structs[s]; ok {
		switch st.Kind {
		case "handler", "object":
			return publicTypeNameForCName(st.CName), true
		case "data":
			return "*" + publicTypeNameForCName(s), true
		}
	}

	return "", false
}

// IsEnumType returns true if the given C type resolves to an enum.
func (r *TypeRegistry) IsEnumType(ctype string) bool {
	ct := strings.TrimSpace(ctype)
	_, ok := r.enums[ct]
	return ok
}

// IsGetterCallback returns true if the field is a function callback with
// exactly one parameter (self) that returns a pointer to a handler-kind struct.
func (r *TypeRegistry) IsGetterCallback(f model.Field) bool {
	if !f.IsFunction {
		return false
	}
	if len(f.Params) != 1 {
		return false
	}
	// Check if the return type is a pointer to a handler-kind struct.
	ret := strings.TrimSpace(f.ReturnCType)
	if ret == "" || ret == "void" {
		return false
	}

	// Try to resolve the return type as a struct pointer.
	s := ret
	s = strings.TrimPrefix(s, "const ")
	if !strings.HasPrefix(s, "struct _") {
		return false
	}
	s = strings.TrimSuffix(s, "*")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "struct ")

	// Look up in registry.
	if st, ok := r.structs[s]; ok {
		return st.Kind == "handler"
	}
	trimmed := strings.TrimLeft(s, "_")
	if st, ok := r.structs[trimmed]; ok {
		return st.Kind == "handler"
	}
	return false
}

// IsHandlerType returns true if the given C type resolves to a handler interface specifically.
func (r *TypeRegistry) IsHandlerType(ctype string) bool {
	if st := r.lookupStructForCType(ctype); st != nil {
		return st.Kind == "handler"
	}
	return false
}

// IsInterfaceType returns true if the given C type resolves to a handler or object interface.
func (r *TypeRegistry) IsInterfaceType(ctype string) bool {
	if st := r.lookupStructForCType(ctype); st != nil {
		return st.Kind == "handler" || st.Kind == "object"
	}
	return false
}

// IsStringType returns true if the C type is a string type (cef_string_userfree_t, cef_string_t*, etc.)
func (r *TypeRegistry) IsStringType(ctype string) bool {
	ct := strings.TrimSpace(ctype)
	switch ct {
	case "cef_string_userfree_t", "const cef_string_t*", "cef_string_t*", "const char*", "char*":
		return true
	}
	return false
}

// IsUserfreeString returns true if the C type is cef_string_userfree_t.
func (r *TypeRegistry) IsUserfreeString(ctype string) bool {
	return strings.TrimSpace(ctype) == "cef_string_userfree_t"
}

// IsDataStructType returns true if the given C type resolves to a data struct pointer.
func (r *TypeRegistry) IsDataStructType(ctype string) bool {
	if st := r.lookupStructForCType(ctype); st != nil {
		return st.Kind == "data"
	}
	return false
}

// lookupStructForCType resolves a C type string to its struct registry entry.
// Handles both "struct _cef_xxx_t*" and bare "cef_xxx_t*" patterns.
func (r *TypeRegistry) lookupStructForCType(ctype string) *model.Struct {
	ct := normalizeConst(strings.TrimSpace(ctype))
	ct = strings.TrimPrefix(ct, "const ")
	if !strings.HasSuffix(ct, "*") {
		return nil
	}
	s := ct
	s = strings.TrimRight(s, "*")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "struct ")

	// Try as-is (e.g., "_cef_browser_t").
	if st, ok := r.structs[s]; ok {
		return st
	}
	// Try without leading underscore (e.g., "cef_browser_t").
	trimmed := strings.TrimLeft(s, "_")
	if st, ok := r.structs[trimmed]; ok {
		return st
	}
	// Try with leading underscore added (bare "cef_xxx_t" → "_cef_xxx_t").
	if !strings.HasPrefix(s, "_") {
		if st, ok := r.structs["_"+s]; ok {
			return st
		}
	}
	return nil
}

// FixupEnumFieldTypes rewrites GoType for enum-typed fields in data structs
// from their internal Go enum type (e.g. CEFStateT) to int32. This is
// necessary because public enums are type-aliased to int32, and the public
// structs are type-aliased to the raw structs — using the internal enum type
// would create a type mismatch for consumers.
func (r *TypeRegistry) FixupEnumFieldTypes(headers []*model.Header) {
	for _, h := range headers {
		for i := range h.Structs {
			s := &h.Structs[i]
			if s.Kind != "data" {
				continue
			}
			for j := range s.Fields {
				f := &s.Fields[j]
				if f.IsFunction {
					continue
				}
				if _, ok := r.enums[f.CType]; ok {
					f.GoType = "int32"
				}
			}
		}
	}
}

// IsBoolReturn returns true if a callback field returns int and its name
// starts with a boolean-like prefix (is_, has_, can_, do_, on_before_).
func IsBoolReturn(f model.Field) bool {
	if !f.IsFunction {
		return false
	}
	ret := strings.TrimSpace(f.ReturnCType)
	if ret != "int" {
		return false
	}
	name := f.CName
	for _, prefix := range []string{"is_", "has_", "can_", "do_", "on_before_"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
