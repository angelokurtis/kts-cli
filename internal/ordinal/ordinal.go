package ordinal

import (
	"errors"
	"strings"
)

func convertNumberToWords(number int) string {
	underTwenty := []string{
		"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten",
		"Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen",
	}
	tens := []string{"Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}
	var words []string

	if number == 0 {
		return "Zero"
	}

	if number >= 1000 {
		thousands := number / 1000
		words = append(words, convertNumberToWords(thousands), "Thousand")
		number %= 1000
	}

	if number >= 100 {
		hundreds := number / 100
		words = append(words, convertNumberToWords(hundreds), "Hundred")
		number %= 100
	}

	if number >= 20 {
		t := number / 10
		words = append(words, tens[t-2])
		number %= 10
	}

	if number > 0 && number < 20 {
		words = append(words, underTwenty[number-1])
	}

	return strings.Join(words, " ")
}

func FromNumber(number int) (string, error) {
	if number < 1 {
		return "", errors.New("input must be a positive integer")
	}

	words := convertNumberToWords(number)

	// Apply ordinal suffix adjustments
	switch {
	case strings.HasSuffix(words, "One"):
		words = strings.TrimSuffix(words, "One") + "First"
	case strings.HasSuffix(words, "Two"):
		words = strings.TrimSuffix(words, "Two") + "Second"
	case strings.HasSuffix(words, "Three"):
		words = strings.TrimSuffix(words, "Three") + "Third"
	case strings.HasSuffix(words, "Five"):
		words = strings.TrimSuffix(words, "Five") + "Fifth"
	case strings.HasSuffix(words, "Eight"):
		words = strings.TrimSuffix(words, "Eight") + "Eighth"
	case strings.HasSuffix(words, "Nine"):
		words = strings.TrimSuffix(words, "Nine") + "Ninth"
	case strings.HasSuffix(words, "Twelve"):
		words = strings.TrimSuffix(words, "Twelve") + "Twelfth"
	default:
		words += "th"
	}

	return words, nil
}
