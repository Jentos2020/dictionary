package controller

import (
	"context"
	dbconverters "leetgo/internal/app/store/converters"
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/entity"
	"leetgo/internal/errors"
)

func (c *Controller) AddWord(ctx context.Context, reqWord entity.Word) error {
	if err := c.DB.AddWord(ctx, dbconverters.EntityToDBWord(reqWord), reqWord.Dictionary); err != nil {
		return err
	}
	currentTrie := c.Trie.Load().(*entity.Trie)
	currentTrie.Insert(reqWord.Data)
	return nil
}

func (c *Controller) RemoveWord(ctx context.Context, word string, dict string) error {
	if word == "" {
		return errors.NewF("empty word")
	}

	if dict != "" {
		removed, err := c.DB.RemoveWord(ctx, dbmodel.Word{Data: word}, dict)
		if err != nil {
			return errors.NewF("delete word from %s: %w", dict, err)
		}
		if !removed {
			return errors.ErrNotFound
		}

		currentTrie := c.Trie.Load().(*entity.Trie)
		currentTrie.Delete(word)

		if err := c.RemoveDefinition(ctx, word, dict); err != nil && !errors.Is(err, errors.ErrNotFound) {
			c.Logger.Warn("Error removing definition during word removal", "error", err, "word", word, "dict", dict)
		} else if errors.Is(err, errors.ErrNotFound) || err == nil {
			c.Logger.Info("Definition removed during word removal", "word", word, "dict", dict)
		}

		return nil
	}

	tables, err := c.DB.GetDictionaryTables(ctx, c.Cfg.Schema)
	if err != nil {
		return errors.NewF("get dictionary tables: %w", err)
	}

	found := false
	for _, t := range tables {
		removed, err := c.DB.RemoveWord(ctx, dbmodel.Word{Data: word}, t)
		if err != nil {
			return errors.NewF("delete word from %s: %w", t, err)
		}
		if removed {
			found = true
			if err := c.RemoveDefinition(ctx, word, t); err != nil && !errors.Is(err, errors.ErrNotFound) {
				c.Logger.Warn("Error removing definition in multi-dict removal", "error", err, "word", word, "dict", t)
			}
		}
	}

	if !found {
		return errors.ErrNotFound
	}

	currentTrie := c.Trie.Load().(*entity.Trie)
	currentTrie.Delete(word)
	return nil
}

func (c *Controller) UpdateWord(ctx context.Context, oldWord, newWord, dict string) error {
	if oldWord == "" || newWord == "" {
		return errors.NewF("old or new word is empty")
	}

	if dict != "" {
		updated, err := c.DB.UpdateWord(ctx, dbmodel.Word{Data: oldWord}, dbmodel.Word{Data: newWord}, dict)
		if err != nil {
			return errors.NewF("update word in %s: %w", dict, err)
		}
		if !updated {
			return errors.ErrNotFound
		}

		currentTrie := c.Trie.Load().(*entity.Trie)
		currentTrie.Delete(oldWord)
		currentTrie.Insert(newWord)

		oldDef, found, err := c.GetDefinition(ctx, oldWord, dict)
		if err != nil {
			c.Logger.Warn("Error checking definition during word update", "error", err, "word", oldWord, "dict", dict)
			return nil
		}
		if found {
			if err := c.RemoveDefinition(ctx, oldWord, dict); err != nil {
				c.Logger.Warn("Error removing old definition during word update", "error", err, "word", oldWord, "dict", dict)
			}
			if err := c.AddDefinition(ctx, newWord, oldDef, dict); err != nil {
				c.Logger.Warn("Error adding new definition during word update", "error", err, "word", newWord, "dict", dict)
			}
			c.Logger.Info("Definition synced during word update", "oldWord", oldWord, "newWord", newWord, "dict", dict)
		}

		return nil
	}

	tables, err := c.DB.GetDictionaryTables(ctx, c.Cfg.Schema)
	if err != nil {
		return errors.NewF("get dictionary tables: %w", err)
	}

	updated := false
	for _, t := range tables {
		ok, err := c.DB.UpdateWord(ctx, dbmodel.Word{Data: oldWord}, dbmodel.Word{Data: newWord}, t)
		if err != nil {
			return errors.NewF("update word in %s: %w", t, err)
		}
		if ok {
			updated = true
			oldDef, found, err := c.GetDefinition(ctx, oldWord, t)
			if err != nil {
				c.Logger.Warn("Error checking definition during multi-dict update", "error", err, "word", oldWord, "dict", t)
				continue
			}
			if found {
				if err := c.RemoveDefinition(ctx, oldWord, t); err != nil {
					c.Logger.Warn("Error removing old definition in multi-dict update", "error", err, "word", oldWord, "dict", t)
				}
				if err := c.AddDefinition(ctx, newWord, oldDef, t); err != nil {
					c.Logger.Warn("Error adding new definition in multi-dict update", "error", err, "word", newWord, "dict", t)
				}
				c.Logger.Info("Definition synced in multi-dict update", "oldWord", oldWord, "newWord", newWord, "dict", t)
			}
		}
	}

	if !updated {
		return errors.ErrNotFound
	}

	currentTrie := c.Trie.Load().(*entity.Trie)
	currentTrie.Delete(oldWord)
	currentTrie.Insert(newWord)
	return nil
}
