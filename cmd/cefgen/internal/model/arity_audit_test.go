package model

import (
	"strings"
	"testing"
)

func TestAuditRawMethodArityAcceptsBoundaryCases(t *testing.T) {
	exactlyFourteen := []Param{{CName: "self", GoName: "self"}}
	for range 14 {
		exactlyFourteen = append(exactlyFourteen, Param{CName: "argument"})
	}

	headers := []*Header{{
		Structs: []Struct{{
			Fields: []Field{
				{CName: "zero", GoName: "Zero", IsFunction: true, Params: []Param{{CName: "self", GoName: "self"}}},
				{CName: "maximum", GoName: "Maximum", IsFunction: true, Params: exactlyFourteen},
			},
		}},
	}}

	if err := AuditRawMethodArity(headers); err != nil {
		t.Fatalf("AuditRawMethodArity() rejected supported arity: %v", err)
	}
}

func TestAuditRawMethodArityAcceptsFourteenArgumentsWithoutSelf(t *testing.T) {
	params := make([]Param, 14)
	headers := []*Header{{Structs: []Struct{{Fields: []Field{{
		CName: "no_self_boundary", GoName: "NoSelfBoundary", IsFunction: true, Params: params,
	}}}}}}

	if err := AuditRawMethodArity(headers); err != nil {
		t.Fatalf("AuditRawMethodArity() rejected 14 arguments without self: %v", err)
	}
}

func TestAuditRawMethodArityRejectsFifteenArgumentsWithoutSelfByCAndGoName(t *testing.T) {
	params := make([]Param, 15)
	headers := []*Header{{Structs: []Struct{{Fields: []Field{{
		CName: "no_self_overflow", GoName: "NoSelfOverflow", IsFunction: true, Params: params,
	}}}}}}

	err := AuditRawMethodArity(headers)
	if err == nil {
		t.Fatal("AuditRawMethodArity() accepted 15 arguments without self")
	}
	for _, want := range []string{"no_self_overflow", "NoSelfOverflow"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("AuditRawMethodArity() error %q does not identify %q", err, want)
		}
	}
}

func TestAuditRawMethodArityReportsEveryMethodByCAndGoName(t *testing.T) {
	tooMany := make([]Param, 16) // self plus 15 arguments
	tooMany[0] = Param{CName: "self", GoName: "self"}
	headers := []*Header{
		{Structs: []Struct{{Fields: []Field{{
			CName: "first_overflow", GoName: "FirstOverflow", IsFunction: true, Params: tooMany,
		}}}}},
		{Structs: []Struct{{Fields: []Field{{
			CName: "second_overflow", GoName: "SecondOverflow", IsFunction: true, Params: tooMany,
		}}}}},
	}

	err := AuditRawMethodArity(headers)
	if err == nil {
		t.Fatal("AuditRawMethodArity() accepted methods with 15 arguments after self")
	}
	for _, want := range []string{"first_overflow", "FirstOverflow", "second_overflow", "SecondOverflow"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("AuditRawMethodArity() error %q does not identify %q", err, want)
		}
	}
}
