package ordinal_test

import (
	"testing"

	"github.com/angelokurtis/kts-cli/internal/ordinal"
)

func TestFromInt(t *testing.T) {
	tests := []struct {
		name   string
		number int
		want   string
	}{
		{"Test for 1", 1, "First"},
		{"Test for 22", 22, "Twenty-Second"},
		{"Test for 45", 45, "Forty-Fifth"},
		{"Test for 100", 100, "One Hundredth"},
		{"Test for 101", 101, "One Hundred First"},
		{"Test for 113", 113, "One Hundred Thirteenth"},
		{"Test for 122", 122, "One Hundred Twenty-Second"},
		{"Test for 1000", 1000, "One Thousandth"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ordinal.FromNumber(tt.number); got != tt.want {
				t.Errorf("ordinal(%d) = %s; want %s", tt.number, got, tt.want)
			}
		})
	}
}
