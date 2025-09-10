// https://leetcode.com/problems/reorder-list/description/

package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestCarFleet(t *testing.T) {
	tests := []struct {
		name     string
		target   int
		position []int
		speed    []int
		expected int
	}{
		{
			name:     "Example 1",
			target:   12,
			position: []int{10, 8, 0, 5, 3},
			speed:    []int{2, 4, 1, 1, 3},
			expected: 3,
		},
		{
			name:     "Example 2",
			target:   10,
			position: []int{3},
			speed:    []int{3},
			expected: 1,
		},
		{
			name:     "Example 3",
			target:   100,
			position: []int{0, 2, 4},
			speed:    []int{4, 2, 1},
			expected: 1,
		},
		{
			name:     "No cars",
			target:   10,
			position: []int{},
			speed:    []int{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := carFleet(tt.target, tt.position, tt.speed)
			if result != tt.expected {
				t.Errorf("Expected %d fleets, got %d", tt.expected, result)
			} else {
				fmt.Printf("%s: Passed, fleets = %d\n", tt.name, result)
			}
		})
	}
}

func carFleet(target int, position []int, speed []int) int {
	n := len(position)
	if n == 0 {
		return 0
	}

	// слайс машин: каждая машина = [позиция, время до цели]
	cars := make([][2]float64, n)
	for i := 0; i < n; i++ {
		cars[i][0] = float64(position[i])
		cars[i][1] = float64(target-position[i]) / float64(speed[i])
	}

	// сортировка машин по позиции по убыванию (от ближней к цели к дальней)
	sort.Slice(cars, func(i, j int) bool {
		return cars[i][0] > cars[j][0]
	})

	fleets := 0
	maxTime := 0.0

	// считаем флоты
	for _, car := range cars {
		time := car[1]
		if time > maxTime {
			fleets++
			maxTime = time
		}
		// ecли time <= maxTime, значит машина догоняет предыдущий флот → флота не увеличиваем
	}

	return fleets
}
