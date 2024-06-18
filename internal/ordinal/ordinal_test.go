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
		{"Test for 1", 1, "first"},
		{"Test for 10", 10, "tenth"},
		{"Test for 11", 11, "eleventh"},
		{"Test for 12", 12, "twelfth"},
		{"Test for 13", 13, "thirteenth"},
		{"Test for 14", 14, "fourteenth"},
		{"Test for 15", 15, "fifteenth"},
		{"Test for 16", 16, "sixteenth"},
		{"Test for 17", 17, "seventeenth"},
		{"Test for 18", 18, "eighteenth"},
		{"Test for 19", 19, "nineteenth"},
		{"Test for 20", 20, "twentieth"},
		{"Test for 30", 30, "thirtieth"},
		{"Test for 40", 40, "fortieth"},
		{"Test for 50", 50, "fiftieth"},
		{"Test for 60", 60, "sixtieth"},
		{"Test for 70", 70, "seventieth"},
		{"Test for 80", 80, "eightieth"},
		{"Test for 90", 90, "ninetieth"},
		{"Test for 22", 22, "twenty-second"},
		{"Test for 45", 45, "forty-fifth"},
		{"Test for 100", 100, "one hundredth"},
		{"Test for 101", 101, "one hundred and first"},
		{"Test for 113", 113, "one hundred and thirteenth"},
		{"Test for 122", 122, "one hundred and twenty-second"},
		{"Test for 148", 148, "one hundred and forty-eighth"},
		{"Test for 1000", 1000, "one thousandth"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ordinal.FromNumber(tt.number); got != tt.want {
				t.Errorf("ordinal(%d) = %s; want %s", tt.number, got, tt.want)
			}
		})
	}
}
