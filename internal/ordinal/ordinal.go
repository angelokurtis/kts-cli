package ordinal

import (
	"errors"
	"strings"

	"github.com/martinusso/inflect"
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func FromNumber[T Number](number T) (string, error) {
	if number < 1 {
		return "", errors.New("input must be a positive integer")
	}

	words := inflect.IntoWords(float64(number))

	// Apply ordinal suffix adjustments
	switch {
	case strings.HasSuffix(words, "one"):
		words = strings.TrimSuffix(words, "one") + "first"
	case strings.HasSuffix(words, "two"):
		words = strings.TrimSuffix(words, "two") + "second"
	case strings.HasSuffix(words, "three"):
		words = strings.TrimSuffix(words, "three") + "third"
	case strings.HasSuffix(words, "five"):
		words = strings.TrimSuffix(words, "five") + "fifth"
	case strings.HasSuffix(words, "eight"):
		words = strings.TrimSuffix(words, "eight") + "eighth"
	case strings.HasSuffix(words, "nine"):
		words = strings.TrimSuffix(words, "nine") + "ninth"
	case strings.HasSuffix(words, "twelve"):
		words = strings.TrimSuffix(words, "twelve") + "twelfth"
	case strings.HasSuffix(words, "ty"):
		words = strings.TrimSuffix(words, "ty") + "tieth"
	default:
		words += "th"
	}

	return words, nil
}
