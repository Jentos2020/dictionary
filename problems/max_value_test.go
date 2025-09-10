// https://leetcode.com/problems/reorder-list/description/

package main

import (
	"strconv"
	"strings"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func TestMaxValue(t *testing.T) {
	tests := []struct {
		name     string
		n        string
		x        int
		expected string
	}{
		{
			name:     "Positive number insert middle",
			n:        "73",
			x:        6,
			expected: "763",
		},
		{
			name:     "Negative number insert beginning",
			n:        "-55",
			x:        2,
			expected: "-255",
		},
		{
			name:     "All digits same positive",
			n:        "99",
			x:        9,
			expected: "999",
		},
		{
			name:     "Negative number insert middle",
			n:        "-13",
			x:        2,
			expected: "-123",
		},
		{
			name:     "Insert at end positive",
			n:        "123",
			x:        1,
			expected: "1231",
		},
		{
			name:     "Insert at end negative",
			n:        "-321",
			x:        5,
			expected: "-3215",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxValue(tt.n, tt.x)
			if got != tt.expected {
				t.Errorf("maxValue(%q, %d) = %q; expected %q", tt.n, tt.x, got, tt.expected)
			}
		})
	}
}

func maxValue(n string, x int) string {
	xChar := strconv.Itoa(x)
	var sb strings.Builder
	inserted := false

	if n[0] == '-' {
		sb.WriteByte('-')
		for i := 1; i < len(n); i++ {
			if !inserted && int(n[i]-'0') > x {
				sb.WriteString(xChar)
				inserted = true
			}
			sb.WriteByte(n[i])
		}
	} else {
		for i := 0; i < len(n); i++ {
			if !inserted && int(n[i]-'0') < x {
				sb.WriteString(xChar)
				inserted = true
			}
			sb.WriteByte(n[i])
		}
	}

	if !inserted {
		sb.WriteString(xChar)
	}

	return sb.String()
}
