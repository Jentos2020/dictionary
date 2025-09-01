// https://leetcode.com/problems/longest-word-in-dictionary/description/

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TrieNode struct {
	IsEnd bool
	Next  map[rune]*TrieNode
}

type Trie struct {
	Root *TrieNode
}

func TestLongestWord(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Example 1",
			input:    []string{"w", "wo", "a", "all", "wor", "worl", "world"},
			expected: "world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := longestWord(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func NewTrie() Trie {
	return Trie{
		Root: &TrieNode{
			Next: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(word string) {
	node := t.Root
	for _, c := range word {
		if _, ok := node.Next[c]; !ok {
			node.Next[c] = &TrieNode{
				Next: make(map[rune]*TrieNode),
			}
		}
		node = node.Next[c]
	}
	node.IsEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.Root
	for _, c := range word {
		if _, ok := node.Next[c]; !ok {
			return false
		}
		node = node.Next[c]
	}
	return node.IsEnd
}

func longestWord(words []string) string {
	trie := NewTrie()
	for _, word := range words {
		trie.Insert(word)
	}

	var longest string
	for _, w := range words {
		if len(w) > len(longest) || len(w) == len(longest) && w < longest {
			if prefixesChecking(trie, w) {
				longest = w
			}
		}
	}

	return longest
}

func prefixesChecking(t Trie, word string) bool {
	for i := 1; i <= len(word); i++ {
		if !t.Search(word[:i]) {
			return false
		}
	}
	return true
}
