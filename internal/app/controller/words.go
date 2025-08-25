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
		deleted, err := c.DB.RemoveWord(ctx, dbmodel.Word{Data: word}, dict)
		if err != nil {
			return errors.NewF("delete word from %s: %w", dict, err)
		}
		if !deleted {
			return errors.ErrNotFound
		}

		currentTrie := c.Trie.Load().(*entity.Trie)
		currentTrie.Delete(word)
		return nil
	}

	tables, err := c.DB.GetDictionaryTables(ctx, c.Cfg.Schema)
	if err != nil {
		return errors.NewF("get dictionary tables: %w", err)
	}

	found := false
	for _, t := range tables {
		deleted, err := c.DB.RemoveWord(ctx, dbmodel.Word{Data: word}, t)
		if err != nil {
			return errors.NewF("delete word from %s: %w", t, err)
		}
		if deleted {
			found = true
			break
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
			break
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
