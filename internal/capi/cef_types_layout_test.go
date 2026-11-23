package capi

import (
	"reflect"
	"testing"
	"unsafe"
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

// TestCEFSettingsTailOffsets pins the offsets of the last cef_settings_t
// fields on 64-bit targets, computed by hand from include/internal/cef_types.h
// (cef_string_t is 24 bytes, 8-byte aligned). A shift here means writes like
// DisableSignalHandlers land in the wrong C field.
func TestCEFSettingsTailOffsets(t *testing.T) {
	if unsafe.Sizeof(uintptr(0)) != 8 {
		t.Skip("offsets computed for 64-bit targets")
	}
	var s CEFSettingsT
	for _, tc := range []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"ChromePolicyID", unsafe.Offsetof(s.ChromePolicyID), 408},
		{"ChromeAppIconID", unsafe.Offsetof(s.ChromeAppIconID), 432},
		{"DisableSignalHandlers", unsafe.Offsetof(s.DisableSignalHandlers), 436},
		{"UseViewsDefaultPopup", unsafe.Offsetof(s.UseViewsDefaultPopup), 440},
		{"sizeof", unsafe.Sizeof(s), 448},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %d, want %d", tc.name, tc.got, tc.want)
		}
	}
}
