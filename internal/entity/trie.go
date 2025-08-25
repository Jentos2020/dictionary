package entity

import (
	"leetgo/internal/consts"
	"sync"
)

type TrieNode struct {
	Children map[rune]*TrieNode
	IsEnd    bool
}

type Trie struct {
	Root *TrieNode
	mu   sync.RWMutex
}

func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[rune]*TrieNode)}}
}

func (t *Trie) Insert(word string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	node := t.Root
	for _, char := range word {
		if _, ok := node.Children[char]; !ok {
			node.Children[char] = &TrieNode{Children: make(map[rune]*TrieNode)}
		}
		node = node.Children[char]
	}
	node.IsEnd = true
}

func (t *Trie) GetWordsByPrefix(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	node := t.Root
	for _, char := range prefix {
		if _, ok := node.Children[char]; !ok {
			return []string{}
		}
		node = node.Children[char]
	}

	words := []string{}
	t.collectWords(node, prefix, &words)
	return words
}

func (t *Trie) collectWords(node *TrieNode, current string, words *[]string) {
	if len(*words) >= consts.TrieSearchLimit {
		return
	}

	if node.IsEnd {
		*words = append(*words, current)
	}

	for char, child := range node.Children {
		if len(*words) >= consts.TrieSearchLimit {
			break
		}
		t.collectWords(child, current+string(char), words)
	}
}

func (t *Trie) Delete(word string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	node := t.Root
	var path []*TrieNode
	path = append(path, node)

	for _, char := range word {
		if child, ok := node.Children[char]; ok {
			node = child
			path = append(path, node)
		} else {
			return // Этого слова и так нету
		}
	}
	if !node.IsEnd {
		return
	}
	node.IsEnd = false

	// Удаляем ноды, если нет потомков
	for i := len(path) - 1; i > 0; i-- {
		parent := path[i-1]
		char := []rune(word)[i-1]
		if len(path[i].Children) == 0 && !path[i].IsEnd { // если в Children лежит пустая мапа и path[i].IsEnd != true
			delete(parent.Children, char)
		} else {
			break
		}
	}
}

func (t *Trie) Copy() *Trie {
	t.mu.RLock()
	defer t.mu.RUnlock()

	newTrie := &Trie{
		Root: t.copyNode(t.Root),
		mu:   sync.RWMutex{},
	}
	return newTrie
}

func (t *Trie) copyNode(node *TrieNode) *TrieNode {
	if node == nil {
		return nil
	}
	newNode := &TrieNode{
		Children: make(map[rune]*TrieNode, len(node.Children)),
		IsEnd:    node.IsEnd,
	}
	for char, child := range node.Children {
		newNode.Children[char] = t.copyNode(child)
	}
	return newNode
}
