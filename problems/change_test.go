package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLemonadeChange(t *testing.T) {
	tests := []struct {
		name     string
		bills    []int
		expected bool
	}{
		{
			name:     "Example 1",
			bills:    []int{5, 5, 5, 10, 20},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lemonadeChange(tt.bills)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func lemonadeChange(bills []int) bool {
	five := 0
	ten := 0

	for _, bill := range bills {
		if bill == 5 {
			five += 1
		} else if bill == 10 {
			if five == 0 {
				return false
			}
			five -= 1
			ten += 1
		} else {
			if ten > 0 && five > 0 {
				five -= 1
				ten -= 1
			} else if five >= 3 {
				five -= 3
			} else {
				return false
			}
		}
	}
	return true
}
