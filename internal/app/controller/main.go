package controller

import (
	"context"
	"leetgo/config"
	"leetgo/internal/app/store"
	"leetgo/internal/app/store/pg"
	"leetgo/internal/entity"
	"log/slog"
	"sync/atomic"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type Controller struct {
	Ctx      context.Context
	Cfg      config.Config
	DB       store.Repository
	Logger   *slog.Logger
	Trie     atomic.Value
	DefCache *simplelru.LRU[string, string]
}

func New(ctx context.Context, cfg config.Config, db *pg.PGStore, log *slog.Logger) *Controller {
	c := &Controller{
		Ctx:    ctx,
		Cfg:    cfg,
		DB:     db,
		Logger: log,
		Trie:   atomic.Value{},
	}
	c.Trie.Store(entity.NewTrie())

	cache, err := simplelru.NewLRU[string, string](3, nil)
	if err != nil {
		log.Error("Failed to create LRU cache", "error", err)
	}
	c.DefCache = cache

	return c
}
