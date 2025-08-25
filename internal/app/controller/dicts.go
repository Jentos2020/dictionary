package controller

import (
	"context"
	"fmt"
	"leetgo/internal/app/controller/converters"
	"leetgo/internal/entity"
)

func (c *Controller) FillTrieWithWords(ctx context.Context, tableName string) error {
	allWords, err := c.DB.GetWords(ctx, tableName) // ~3-7 сек
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Failed to get words from DB: %v", err))
		return err
	}

	wordsEntity := converters.DBWordsToEntity(allWords)

	newTrie := entity.NewTrie()
	for _, w := range wordsEntity { // ~1.4 сек
		newTrie.Insert(w.Data)
	}

	c.Trie.Store(newTrie)
	c.Logger.Info("Trie loaded")

	return nil
}
