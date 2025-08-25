package controller

import (
	"context"
	"leetgo/config"
	"leetgo/internal/app/store"
	"leetgo/internal/app/store/pg"
	"leetgo/internal/entity"
	"log/slog"
	"sync/atomic"
)

type Controller struct {
	Ctx    context.Context
	Cfg    config.Config
	DB     store.Repository
	Logger *slog.Logger
	Trie   atomic.Value
}

func New(
	ctx context.Context,
	cfg config.Config,
	db *pg.PGStore,
	log *slog.Logger) *Controller {
	c := &Controller{
		Ctx:    ctx,
		Cfg:    cfg,
		DB:     db,
		Logger: log,
	}
	c.Trie.Store(entity.NewTrie())
	return c
}
