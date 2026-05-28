package tags

import "testing"

func Test_convertToSemVer(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{name: "empty", in: "", want: "", ok: false},
		{name: "no version", in: "latest", want: "", ok: false},
		{name: "major", in: "1", want: "1.0.0", ok: true},
		{name: "major minor", in: "1.2", want: "1.2.0", ok: true},
		{name: "major minor patch", in: "1.2.3", want: "1.2.3", ok: true},
		{name: "prefixed version", in: "v1.2.3", want: "1.2.3", ok: true},
		{name: "version in tag", in: "release-1.2.3-alpine", want: "1.2.3", ok: true},
		{name: "truncates extra parts", in: "1.2.3.4", want: "1.2.3", ok: true},
		{name: "first numeric version wins", in: "15beta3-alpine3.16", want: "15.0.0", ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := convertToSemVer(tt.in)

			if got != tt.want || ok != tt.ok {
				t.Fatalf("convertToSemVer(%q) = (%q, %v), want (%q, %v)", tt.in, got, ok, tt.want, tt.ok)
			}
		})
	}
}
