package parser

import "testing"

func TestPreprocessCEFAPIBoundaries(t *testing.T) {
	tests := []struct {
		name   string
		target int
		input  string
		want   string
	}{
		{
			name:   "added before target",
			target: 15000,
			input:  "#if CEF_API_ADDED(14999)\nselected\n#else\nrejected\n#endif",
			want:   "selected",
		},
		{
			name:   "added at target",
			target: 15000,
			input:  "#if CEF_API_ADDED(15000)\nselected\n#else\nrejected\n#endif",
			want:   "selected",
		},
		{
			name:   "removed before target",
			target: 15000,
			input:  "#if CEF_API_REMOVED(14999)\nrejected\n#else\nselected\n#endif",
			want:   "selected",
		},
		{
			name:   "removed at target",
			target: 15000,
			input:  "#if CEF_API_REMOVED(15000)\nrejected\n#else\nselected\n#endif",
			want:   "selected",
		},
		{
			name:   "range lower boundary is inclusive",
			target: 15000,
			input:  "#if CEF_API_RANGE(15000, 16000)\nselected\n#else\nrejected\n#endif",
			want:   "selected",
		},
		{
			name:   "range upper boundary is exclusive",
			target: 16000,
			input:  "#if CEF_API_RANGE(15000, 16000)\nrejected\n#else\nselected\n#endif",
			want:   "selected",
		},
		{
			name:   "nested selection",
			target: 15000,
			input:  "#if CEF_API_ADDED(14000)\nouter\n#if CEF_API_RANGE(15000, 16000)\nnested\n#else\nrejected nested\n#endif\n#else\nrejected outer\n#endif",
			want:   "outer\nnested",
		},
		{
			name:   "experimental remains future",
			target: 15000,
			input:  "#if CEF_API_ADDED(CEF_EXPERIMENTAL)\nrejected\n#elif CEF_API_REMOVED(CEF_EXPERIMENTAL)\nselected\n#endif",
			want:   "selected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PreprocessCEFAPI([]byte(tt.input), tt.target)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPreprocessCEFAPIRejectsInvalidExpressions(t *testing.T) {
	for _, input := range []string{
		"#if CEF_API_ADDED(foo)\n#endif",
		"#if CEF_API_RANGE(15000)\n#endif",
		"#if CEF_API_UNKNOWN(15000)\n#endif",
		"#elif CEF_API_ADDED(15000)\n#endif",
		"#if CEF_API_ADDED(15000)\n",
	} {
		if _, err := PreprocessCEFAPI([]byte(input), 15000); err == nil {
			t.Errorf("PreprocessCEFAPI(%q) succeeded", input)
		}
	}
}
