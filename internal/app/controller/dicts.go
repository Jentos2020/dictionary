package controller

import (
	"context"
	"fmt"
	"leetgo/internal/entity"
	"leetgo/internal/errors"
	"runtime"
	"time"
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

	startDB := time.Now()
	allWords, err := c.DB.GetWords(ctx, tableName)
	elapsedDB := time.Since(startDB)

	if err != nil {
		c.Logger.Error(fmt.Sprintf("Failed to get words from DB: %v", err))
		return err
	}

	c.Logger.Info(fmt.Sprintf("Loaded %d words from DB in %s", len(allWords), elapsedDB))

	currentTrie := c.Trie.Load().(*entity.Trie)

	startTotal := time.Now()

	newTrie := currentTrie.Copy()

	var mBefore runtime.MemStats
	runtime.ReadMemStats(&mBefore)

	startInsert := time.Now()
	for _, w := range allWords {
		newTrie.Insert(w.Data)
	}

	var mAfter runtime.MemStats
	runtime.ReadMemStats(&mAfter)
	elapsedInsert := time.Since(startInsert)

	c.Trie.Store(newTrie)

	elapsedTotal := time.Since(startTotal)

	usedBytes := mAfter.Alloc - mBefore.Alloc
	c.Logger.Info(fmt.Sprintf("Approximate Trie size: %.2f MB\n", float64(usedBytes)/1024/1024))
	c.Logger.Info(fmt.Sprintf("Trie built in %s (words: %d)", elapsedTotal, len(allWords)))
	c.Logger.Info(fmt.Sprintf("  └─ Insert: %s", elapsedInsert))

	return nil
}
