package controller

import (
	"context"
	"fmt"
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/consts"
	"leetgo/internal/errors"
)

func (c *Controller) GetDefinition(ctx context.Context, word, dict string) (string, bool, error) {
	if dict == "" {
		dict = consts.DefaultDict
	}
	if word == "" {
		return "", false, errors.NewF("empty word")
	}

	key := fmt.Sprintf("%s:%s", dict, word)
	if c.DefCache != nil {
		if val, ok := c.DefCache.Get(key); ok {
			c.Logger.Debug("Definition cache hit", "key", key)
			return val, true, nil
		}
	}

	def, found, err := c.DB.GetDefinition(ctx, word, dict)
	if err != nil {
		return "", false, errors.NewF("get definition: %w", err)
	}
	if found && c.DefCache != nil {
		c.DefCache.Add(key, def)
		c.Logger.Debug("Definition cache miss -> added", "key", key)
	}
	return def, found, nil
}

func (c *Controller) AddDefinition(ctx context.Context, word, defText, dict string) error {
	if dict == "" {
		dict = consts.DefaultDict
	}
	if word == "" || defText == "" {
		return errors.NewF("empty word or definition")
	}

	if err := c.DB.AddDefinition(ctx, dbmodel.Definition{Word: word, Definition: defText}, dict); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", dict, word)
	if c.DefCache != nil {
		c.DefCache.Remove(key)
	}
	c.Logger.Info("Definition added and cache invalidated", "word", word, "dict", dict)
	return nil
}

func (c *Controller) UpdateDefinition(ctx context.Context, word, newDef, dict string) error {
	if dict == "" {
		dict = consts.DefaultDict
	}
	if word == "" || newDef == "" {
		return errors.NewF("empty word or new definition")
	}

	updated, err := c.DB.UpdateDefinition(ctx, word, newDef, dict)
	if err != nil {
		return err
	}
	if !updated {
		return errors.ErrNotFound
	}

	key := fmt.Sprintf("%s:%s", dict, word)
	if c.DefCache != nil {
		c.DefCache.Remove(key)
	}
	c.Logger.Info("Definition updated and cache invalidated", "word", word, "dict", dict)
	return nil
}

func (c *Controller) RemoveDefinition(ctx context.Context, word, dict string) error {
	if dict == "" {
		dict = consts.DefaultDict
	}
	if word == "" {
		return errors.NewF("empty word")
	}

	removed, err := c.DB.RemoveDefinition(ctx, word, dict)
	if err != nil {
		return err
	}
	if !removed {
		return errors.ErrNotFound
	}

	key := fmt.Sprintf("%s:%s", dict, word)
	if c.DefCache != nil {
		c.DefCache.Remove(key)
	}
	c.Logger.Info("Definition removed and cache invalidated", "word", word, "dict", dict)
	return nil
}
