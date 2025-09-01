// https://leetcode.com/problems/unique-binary-search-trees-ii/description/

package main

import (
	"fmt"
	"testing"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func TestUniqueBinarySearchTrees(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected [][]interface{}
	}{
		{
			name:  "Example n=3",
			input: 3,
			expected: [][]interface{}{
				{1, nil, 2, nil, 3},
				{1, nil, 3, 2},
				{2, 1, 3},
				{3, 1, nil, nil, 2},
				{3, 2, nil, 1},
			},
		},
		{
			name:     "n=1",
			input:    1,
			expected: [][]interface{}{{1}},
		},
		{
			name:     "n=0",
			input:    0,
			expected: [][]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trees := generateTrees(tt.input)
			var results [][]interface{}
			for _, tree := range trees {
				results = append(results, treeToArray(tree))
			}

			// Сравниваем без учёта порядка
			if len(results) != len(tt.expected) {
				t.Errorf("Expected %d trees, got %d", len(tt.expected), len(results))
			}

			for i, arr := range results {
				fmt.Printf("Tree %d: %v\n", i+1, arr)
			}
		})
	}
}

// =========================
// Основная функция
// =========================
func generateTrees(n int) []*TreeNode {
	if n == 0 {
		return []*TreeNode{}
	}
	return generate(1, n)
}

func generate(start, end int) []*TreeNode {
	var allTrees []*TreeNode

	if start > end {
		allTrees = append(allTrees, nil)
		return allTrees
	}

	for i := start; i <= end; i++ {
		leftTrees := generate(start, i-1)
		rightTrees := generate(i+1, end)

		for _, left := range leftTrees {
			for _, right := range rightTrees {
				root := &TreeNode{Val: i, Left: left, Right: right}
				allTrees = append(allTrees, root)
			}
		}
	}

	return allTrees
}

// =========================
// Преобразование дерева в массив
// =========================
func treeToArray(root *TreeNode) []interface{} {
	if root == nil {
		return []interface{}{}
	}
	var result []interface{}
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node != nil {
			result = append(result, node.Val)
			queue = append(queue, node.Left)
			queue = append(queue, node.Right)
		} else {
			result = append(result, nil)
		}
	}

	// Убираем лишние nil в конце
	for len(result) > 0 && result[len(result)-1] == nil {
		result = result[:len(result)-1]
	}

	return result
}
