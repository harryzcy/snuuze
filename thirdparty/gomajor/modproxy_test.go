package gomajor

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		old, new     string
		major, newer bool
	}{
		{
			old:   "v1.0.0",
			new:   "v1.0.1",
			newer: true,
		},
		{
			old:   "v1.0.1",
			new:   "v1.0.0",
			newer: false,
		},
		{
			old:   "v1.0.0",
			new:   "v2.0.0",
			major: true,
			newer: true,
		},
		{
			old:   "v1.0.9",
			new:   "v1.4.0",
			major: true,
			newer: false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if ok := IsNewerVersion(tt.old, tt.new, tt.major); tt.newer != ok {
				t.Fatalf("IsNewerVersion(%q, %q, %v) = %v", tt.old, tt.new, tt.major, ok)
			}
		})
	}
}

func TestCompareVersion(t *testing.T) {
	tests := []struct {
		v, w string
		want int
	}{
		{v: "v0.0.0", w: "v1.0.0", want: -1},
		{v: "v1.0.0", w: "v0.0.0", want: 1},
		{v: "v0.0.0", w: "v0.0.0", want: 0},
		{v: "v12.0.0+incompatible", w: "v0.0.0", want: -1},
		{v: "", w: "", want: 0},
		{v: "v0.1.0", w: "bad", want: 1},
		{v: "v0.0.0+incompatible", w: "v0.0.0", want: -1},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := CompareVersion(tt.v, tt.w); got != tt.want {
				t.Fatalf("CompareVersion(%q, %q) = %v, want %v", tt.v, tt.w, got, tt.want)
			}
		})
	}
}

func TestVersionRange(t *testing.T) {
	tests := []struct {
		r    VersionRange
		v    string
		want bool
	}{
		{
			r:    VersionRange{Low: "v0.0.0", High: "v0.0.1"},
			v:    "v0.0.0",
			want: true,
		},
		{
			r:    VersionRange{Low: "v0.0.0", High: "v0.0.1"},
			v:    "v0.0.1",
			want: true,
		},
		{
			r:    VersionRange{Low: "v0.0.0", High: "v0.0.1"},
			v:    "v0.0.2",
			want: false,
		},
		{
			r:    VersionRange{Low: "v0.0.0", High: "v0.0.1"},
			v:    "v0.0.0+incompatible",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.r.Includes(tt.v); got != tt.want {
				t.Fatalf("VersionRange{Low: %q, High: %q}.Includes(%q) = %v, want %v", tt.r.Low, tt.r.High, tt.v, got, tt.want)
			}
		})
	}
}
