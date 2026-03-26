package model

import "testing"

func TestPublicName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"cef_browser_t", "Browser"},
		{"cef_life_span_handler_t", "LifeSpanHandler"},
		{"_cef_browser_t", "Browser"},
		{"cef_request_context_t", "RequestContext"},
		{"cef_urlrequest_t", "Urlrequest"},
		{"cef_browser_host_t", "BrowserHost"},
		{"cef_command_line_t", "CommandLine"},
		{"cef_url_parts_t", "URLParts"},
		{"cef_process_id_t", "ProcessID"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := PublicName(tt.input)
			if got != tt.want {
				t.Errorf("PublicName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
