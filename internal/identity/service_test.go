package identity

import (
	"errors"
	"testing"
)

func TestNormalizeNewUser(t *testing.T) {
	tests := []struct {
		name    string
		in      NewUser
		want    NewUser
		wantErr error
	}{
		{"trims and lowercases", NewUser{"  Ana ", " Ana@Mail.COM "}, NewUser{"Ana", "ana@mail.com"}, nil},
		{"empty name", NewUser{"  ", "ana@mail.com"}, NewUser{}, ErrInvalidInput},
		{"bad email", NewUser{"Ana", "not-an-email"}, NewUser{}, ErrInvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeNewUser(tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestClampLimit(t *testing.T) {
	cases := map[uint64]uint64{0: defaultListLimit, 10: 10, 10_000: maxListLimit}
	for in, want := range cases {
		if got := clampLimit(in); got != want {
			t.Errorf("clampLimit(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestIsUUID(t *testing.T) {
	cases := map[string]bool{
		"0195f3a2-7c1e-7b3a-9f1e-2a4b6c8d0e1f": true,
		"0195F3A2-7C1E-7B3A-9F1E-2A4B6C8D0E1F": true,
		"not-a-uuid":                           false,
		"0195f3a2x7c1e-7b3a-9f1e-2a4b6c8d0e1f": false,
		"0195f3a2-7c1e-7b3a-9f1e-2a4b6c8d0e1g": false,
	}
	for in, want := range cases {
		if got := isUUID(in); got != want {
			t.Errorf("isUUID(%q) = %v, want %v", in, got, want)
		}
	}
}
