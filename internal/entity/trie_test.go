package entity

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrieInsertAndGetWordsByPrefix(t *testing.T) {
	trie := NewTrie()
	words := []string{"apple", "ant", "art", "banana", "band"}

	// Insert
	for _, w := range words {
		trie.Insert(w)
	}

	// Search by prefix
	result := trie.GetWordsByPrefix("a")
	expected := []string{"ant", "apple", "art"}
	sort.Strings(result)
	sort.Strings(expected)
	assert.Equal(t, expected, result, "Mismatch for prefix 'a'")

	// Empty prefix
	result = trie.GetWordsByPrefix("")
	expectedAll := words
	sort.Strings(result)
	sort.Strings(expectedAll)
	assert.Equal(t, expectedAll, result, "Mismatch for empty prefix")

	// No match
	result = trie.GetWordsByPrefix("z")
	assert.Empty(t, result, "Expected empty for 'z'")
}

func TestTrieDelete(t *testing.T) {
	trie := NewTrie()
	trie.Insert("apple")
	trie.Insert("app")
	trie.Insert("banana")

	// Delete
	trie.Delete("apple")

	result := trie.GetWordsByPrefix("a")
	expected := []string{"app"}
	sort.Strings(result)
	sort.Strings(expected)
	assert.Equal(t, expected, result, "Mismatch after delete")
}

func TestTrieUpdate(t *testing.T) {
	// Update = Delete + Insert
	trie := NewTrie()
	trie.Insert("cast")
	trie.Insert("cama")
	trie.Insert("cat")
	trie.Insert("car")
	trie.Insert("cargo")
	trie.Delete("cast")
	trie.Insert("barco")

	result := trie.GetWordsByPrefix("b")
	expected := []string{"barco"}
	assert.Equal(t, expected, result, "Mismatch after update")
}
