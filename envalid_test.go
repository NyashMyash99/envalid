package envalid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"should return string", "string", "string"},
		{"should return string when input contains one letter", "s", "s"},
		{"should return empty string when input is empty", "", ""},
		{"should return empty string when input is whitespace", " ", ""},
		{"should trim whitespace at beginning", " s", "s"},
		{"should trim whitespace at ending", "s ", "s"},
		{"should trim whitespaces", " s ", "s"},
		{"should trim tabs and newlines", "\ts\n", "s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vRes, vErr := Str(tt.input)

			assert.Equal(t, tt.expected, vRes)
			assert.NoError(t, vErr)
		})
	}
}
