package updater

import (
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		latest   string
		current  string
		expected bool
	}{
		{"v1.0.5", "v1.0.4", true},
		{"v1.1.0", "v1.0.4", true},
		{"v2.0.0", "v1.9.9", true},
		{"v1.0.4", "v1.0.4", false},
		{"v1.0.3", "v1.0.4", false},
		{"v1.0.4.1", "v1.0.4", true},
		{"1.0.5", "1.0.4", true},
		{"", "v1.0.4", false},
		{"v1.0.4", "", false},
	}

	for _, tt := range tests {
		got := IsNewerVersion(tt.latest, tt.current)
		if got != tt.expected {
			t.Errorf("IsNewerVersion(%q, %q) = %v, expected %v", tt.latest, tt.current, got, tt.expected)
		}
	}
}
