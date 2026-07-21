package capi

import (
	"reflect"
	"testing"
)

// These fields are ABI-sensitive CEF API 150 selections. Keeping this test
// outside generated output catches accidental preprocessing regressions.
func TestCEF150SettingsLayouts(t *testing.T) {
	settings := reflect.TypeFor[CEFSettingsT]()
	if _, ok := settings.FieldByName("UseViewsDefaultPopup"); !ok {
		t.Fatal("CEFSettingsT is missing API 150 UseViewsDefaultPopup")
	}
	browser := reflect.TypeFor[CEFBrowserSettingsT]()
	if _, ok := browser.FieldByName("DatabasesDeprecated"); !ok {
		t.Fatal("CEFBrowserSettingsT is missing DatabasesDeprecated")
	}
	if _, ok := browser.FieldByName("Databases"); ok {
		t.Fatal("CEFBrowserSettingsT retained removed Databases field")
	}
}
