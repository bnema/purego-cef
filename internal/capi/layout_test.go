package capi

import (
	"testing"
	"unsafe"
)

// TestStructSizes verifies that generated Go struct sizes match the C ABI
// on linux/amd64. This catches layout bugs (wrong types, missing padding)
// without needing the CEF runtime.
func TestStructSizes(t *testing.T) {
	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		// Primitives / small types
		{"CEFStringT", unsafe.Sizeof(CEFStringT{}), 24},
		{"CEFBasetimeT", unsafe.Sizeof(CEFBasetimeT{}), 8},

		// Geometry types (cef_types_geometry.h)
		{"CEFPointT", unsafe.Sizeof(CEFPointT{}), 8},
		{"CEFRectT", unsafe.Sizeof(CEFRectT{}), 16},
		{"CEFSizeT", unsafe.Sizeof(CEFSizeT{}), 8},
		{"CEFInsetsT", unsafe.Sizeof(CEFInsetsT{}), 16},
		{"CEFRangeT", unsafe.Sizeof(CEFRangeT{}), 8},

		// Input event types
		{"CEFMouseEventT", unsafe.Sizeof(CEFMouseEventT{}), 12},
		{"CEFTouchEventT", unsafe.Sizeof(CEFTouchEventT{}), 40},
		{"CEFKeyEventT", unsafe.Sizeof(CEFKeyEventT{}), 40},

		// Medium structs (cef_types.h)
		{"CEFDraggableRegionT", unsafe.Sizeof(CEFDraggableRegionT{}), 20},
		{"CEFPopupFeaturesT", unsafe.Sizeof(CEFPopupFeaturesT{}), 48},
		{"CEFScreenInfoT", unsafe.Sizeof(CEFScreenInfoT{}), 56},
		{"CEFCursorInfoT", unsafe.Sizeof(CEFCursorInfoT{}), 32},
		{"CEFAudioParametersT", unsafe.Sizeof(CEFAudioParametersT{}), 24},
		{"CEFCompositionUnderlineT", unsafe.Sizeof(CEFCompositionUnderlineT{}), 32},
		{"CEFMediaSinkDeviceInfoT", unsafe.Sizeof(CEFMediaSinkDeviceInfoT{}), 64},
		{"CEFTouchHandleStateT", unsafe.Sizeof(CEFTouchHandleStateT{}), 48},

		// Large structs (cef_types.h)
		{"CEFCookieT", unsafe.Sizeof(CEFCookieT{}), 152},
		{"CEFBoxLayoutSettingsT", unsafe.Sizeof(CEFBoxLayoutSettingsT{}), 56},
		{"CEFPdfPrintSettingsT", unsafe.Sizeof(CEFPdfPrintSettingsT{}), 168},
		{"CEFUrlpartsT", unsafe.Sizeof(CEFUrlpartsT{}), 248},
		{"CEFTaskInfoT", unsafe.Sizeof(CEFTaskInfoT{}), 88},
		{"CEFSettingsT", unsafe.Sizeof(CEFSettingsT{}), 440},
		{"CEFBrowserSettingsT", unsafe.Sizeof(CEFBrowserSettingsT{}), 264},
		{"CEFRequestContextSettingsT", unsafe.Sizeof(CEFRequestContextSettingsT{}), 96},

		// Linux-specific types (cef_types_linux.h)
		{"CEFMainArgsT", unsafe.Sizeof(CEFMainArgsT{}), 16},
		{"CEFWindowInfoT", unsafe.Sizeof(CEFWindowInfoT{}), 88},
		{"CEFLinuxWindowPropertiesT", unsafe.Sizeof(CEFLinuxWindowPropertiesT{}), 104},
		{"CEFAcceleratedPaintNativePixmapPlaneT", unsafe.Sizeof(CEFAcceleratedPaintNativePixmapPlaneT{}), 32},
		{"CEFAcceleratedPaintInfoT", unsafe.Sizeof(CEFAcceleratedPaintInfoT{}), 272},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("sizeof(%s) = %d, want %d", tt.name, tt.got, tt.want)
			}
		})
	}
}
