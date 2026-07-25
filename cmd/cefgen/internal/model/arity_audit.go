package model

import (
	"fmt"
	"strings"
)

const maxRawMethodArgumentsAfterSelf = 14

// RawMethodParams returns the arguments emitted after the implicit receiver.
// CEF callbacks conventionally name that receiver self; functions without it
// retain their first argument.
func RawMethodParams(params []Param) []Param {
	if len(params) == 0 {
		return nil
	}
	if strings.EqualFold(params[0].CName, "self") || strings.EqualFold(params[0].GoName, "self") {
		return params[1:]
	}
	return params
}

// AuditRawMethodArity verifies that every CEF method can be called through the
// fixed-arity SyscallSelf path used by generated raw wrappers.
func AuditRawMethodArity(headers []*Header) error {
	var violations []string
	for _, header := range headers {
		if header == nil {
			continue
		}
		for _, cefStruct := range header.Structs {
			for _, method := range cefStruct.Fields {
				if !method.IsFunction {
					continue
				}
				argumentCount := len(RawMethodParams(method.Params))
				if argumentCount > maxRawMethodArgumentsAfterSelf {
					violations = append(violations, fmt.Sprintf(
						"C method %q / Go method %q has %d arguments after self (maximum %d)",
						method.CName, method.GoName, argumentCount, maxRawMethodArgumentsAfterSelf,
					))
				}
			}
		}
	}
	if len(violations) > 0 {
		return fmt.Errorf("raw method arity audit failed:\n%s", strings.Join(violations, "\n"))
	}
	return nil
}
