// https://leetcode.com/problems/check-if-array-pairs-are-divisible-by-k/description/

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanArrange(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		k        int
		expected bool
	}{
		{
			name:     "Example 1: Valid pairs divisible by k=5",
			input:    []int{1, 2, 3, 4, 5, 10, 6, 7, 8, 9},
			k:        5,
			expected: true,
			// Пары: (1,9), (2,8), (3,7), (4,6), (5,10) — суммы кратны 5.
		},
		{
			name:     "Example 2: Valid pairs divisible by k=7",
			input:    []int{1, 2, 3, 4, 5, 6},
			k:        7,
			expected: true,
			// Пары: (1,6), (2,5), (3,4) — суммы кратны 7.
		},
		{
			name:     "Example 3: No valid pairs for k=10",
			input:    []int{1, 2, 3, 4, 5, 6},
			k:        10,
			expected: false,
			// Нельзя составить 3 пары с суммами, кратными 10.
		},
		{
			name:     "Empty array",
			input:    []int{},
			k:        5,
			expected: true,
			// Пустой массив можно считать разбиваемым на 0 пар, что удовлетворяет условию.
		},
		{
			name:     "Single element",
			input:    []int{1},
			k:        5,
			expected: false,
			// Нельзя составить пары из одного элемента.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canArrange(tt.input, tt.k)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func canArrange(arr []int, k int) bool {
	freq := make(map[int]int)

	for _, n := range arr {
		rem := (n%k + k) % k // чтобы всегда было положительное значение
		freq[rem]++
	}

	if freq[0]%2 == 1 { // если нечетное количество чисел делилось без остатка на k,
		return false // значит не у всех чисел есть нужна пара которая даст в сумме деления на k без остатка
	}

	for rem, _ := range freq {
		if rem == 0 {
			continue
		}
		if freq[rem] != freq[k-rem] {
			return false
		}
	}

	for rem := 1; rem < k; rem++ {
		if freq[rem] != freq[k-rem] {
			return false
		}
	}

	return true
}
