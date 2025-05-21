// Copyright 2023 Intrinsic Innovation LLC

package validate

import (
	"errors"
	"strings"
	"testing"

	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
)

func TestUserString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:  "valid",
			input: "hi_i_am_a_valid_string!",
		},
		{
			name:  "url",
			input: "https://www.intrinsic.ai",
		},
		{
			name:  "empty",
			input: "",
		},
		{
			name: "invalid character",
			input: `
			This has new lines!
			`,
			wantErr: errDoesNotMatchPattern,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := UserString(tc.input); !errors.Is(err, tc.wantErr) {
				t.Errorf("UserString(%q) = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestAlphabetic(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:  "valid",
			input: "validstring",
		},
		{
			name:  "empty",
			input: "",
		},
		{
			name:    "has space",
			input:   "invalid string",
			wantErr: errDoesNotMatchPattern,
		},
		{
			name:    "has number",
			input:   "invalid123string",
			wantErr: errDoesNotMatchPattern,
		},
		{
			name:    "has hyphen",
			input:   "invalid-string",
			wantErr: errDoesNotMatchPattern,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := Alphabetic(tc.input); !errors.Is(err, tc.wantErr) {
				t.Errorf("Alphabetic(%q) = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestDNSLabel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:  "valid",
			input: "valid-label",
		},
		{
			name:  "almost too long",
			input: strings.Repeat("a", dnsLabelMaxLength),
		},
		{
			name:    "empty",
			input:   "",
			wantErr: errDoesNotMatchPattern,
		},
		{
			name:    "too long",
			input:   strings.Repeat("a", dnsLabelMaxLength+1),
			wantErr: errTooLong,
		},
		{
			name:    "invalid character",
			input:   "invalid_label",
			wantErr: errDoesNotMatchPattern,
		},
		{
			name:    "starts with number",
			input:   "1-label",
			wantErr: errDoesNotMatchPattern,
		},
		{
			name:    "starts with dash",
			input:   "-label",
			wantErr: errDoesNotMatchPattern,
		},
		{
			name:    "ends with dash",
			input:   "label-",
			wantErr: errDoesNotMatchPattern,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := DNSLabel(tc.input); !errors.Is(err, tc.wantErr) {
				t.Errorf("DNSLabel(%q) = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestImage(t *testing.T) {
	tests := []struct {
		name    string
		img     *ipb.Image
		wantErr error
	}{
		{
			name: "valid with tag",
			img: &ipb.Image{
				Registry: "gcr.io/test-project",
				Name:     "test-image",
				Tag:      ":latest",
			},
		},
		{
			name: "valid with sha256",
			img: &ipb.Image{
				Registry: "gcr.io/test-project",
				Name:     "test-image",
				Tag:      "@sha256:abcdef0123456789",
			},
		},
		{
			name: "invalid registry",
			img: &ipb.Image{
				Registry: "-gcr.io/test-project",
				Name:     "test-image",
				Tag:      ":latest",
			},
			wantErr: errDoesNotMatchPattern,
		},
		{
			name: "invalid name",
			img: &ipb.Image{
				Registry: "gcr.io/test-project",
				Name:     "-test-image",
				Tag:      ":latest",
			},
			wantErr: errDoesNotMatchPattern,
		},
		{
			name: "invalid tag",
			img: &ipb.Image{
				Registry: "gcr.io/test-project",
				Name:     "test-image",
				Tag:      "invalid-tag",
			},
			wantErr: errDoesNotMatchPattern,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := Image(tc.img); !errors.Is(err, tc.wantErr) {
				t.Errorf("Image(%v) = %v, want %v", tc.img, err, tc.wantErr)
			}
		})
	}
}
