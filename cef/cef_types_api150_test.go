package cef

import "testing"

func TestCEF150ConditionalEnumValues(t *testing.T) {
	if ResultcodeResultCodeInvalidIsolatedBrowserProcess != 40 || ResultcodeResultCodeChromeLast != 41 {
		t.Fatalf("API 150 result-code values = %d, %d", ResultcodeResultCodeInvalidIsolatedBrowserProcess, ResultcodeResultCodeChromeLast)
	}
	if WindowOpenDispositionWodNewSplitView != 12 || WindowOpenDispositionWodNumValues != 13 {
		t.Fatalf("API 150 window-disposition values = %d, %d", WindowOpenDispositionWodNewSplitView, WindowOpenDispositionWodNumValues)
	}
}
