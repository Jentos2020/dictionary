package controller

import (
	"context"
	"fmt"
	"leetgo/internal/entity"
	"leetgo/internal/errors"
)

func (c *Controller) FillTrieWithWords(ctx context.Context, tableName string) error {
	tables, err := c.DB.GetDictionaryTables(ctx, c.Cfg.Schema)
	if err != nil {
		return errors.NewF("get dictionary tables: %w", err)
	}

	var found bool
	for _, table := range tables {
		if table == tableName {
			found = true
			break
		}
	}
	if !found {
		err := c.DB.WriteDictsToDb(ctx, c.Cfg, tableName)
		if err != nil {
			return errors.Wrapf(err, "while creating new table %s", tableName)
		}
	}

	allWords, err := c.DB.GetWords(ctx, tableName)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Failed to get words from DB: %v", err))
		return err
	}

	currentTrie := c.Trie.Load().(*entity.Trie)
	newTrie := currentTrie.Copy()
	for _, w := range allWords {
		newTrie.Insert(w.Data)
	}

	c.Trie.Store(newTrie)
	c.Logger.Info("Trie loaded")

	return nil
}
