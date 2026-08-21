package preferences

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeLabel(t *testing.T) {
	tooLong := strings.Repeat("x", MaxLabelLen+1)
	label := "  Home farm  "
	empty := "   "

	tests := []struct {
		name    string
		input   *string
		want    *string
		wantErr error
	}{
		{name: "omitted", input: nil, want: nil},
		{name: "trimmed", input: &label, want: stringPtr("Home farm")},
		{name: "blank becomes absent", input: &empty, want: nil},
		{name: "too long", input: &tooLong, wantErr: ErrInvalidData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeLabel(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if got == nil || tt.want == nil {
				if got != nil || tt.want != nil {
					t.Fatalf("label = %v, want %v", got, tt.want)
				}
				return
			}
			if *got != *tt.want {
				t.Errorf("label = %q, want %q", *got, *tt.want)
			}
		})
	}
}

func stringPtr(value string) *string { return &value }
