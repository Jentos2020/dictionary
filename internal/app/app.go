package app

import (
	"context"
	"fmt"
	"leetgo/config"
	"leetgo/internal/app/controller"
	"leetgo/internal/app/handler"
	"leetgo/internal/app/store"
	"leetgo/internal/app/store/pg"

	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	cfg  config.Config
	log  *slog.Logger
	seed bool
}

func New(cfg config.Config, log *slog.Logger, seed bool) *App {
	return &App{
		cfg:  cfg,
		log:  log,
		seed: seed,
	}
}

func (app *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := app.cfg.GetDsn()
	db, err := pg.New(
		dsn,
		app.cfg.Schema,
		app.cfg.MaxConn,
		app.cfg.MaxIdle,
		app.log,
	)
	if err != nil {
		return err
	}

	c := controller.New(ctx, app.cfg, db, app.log)
	h := handler.New(c)

	err = store.MigrateUp(ctx, app.cfg, app.log)
	if err != nil {
		return err
	}

	if app.seed {
		err = db.WriteDictsToDb(ctx, app.cfg, "")
		if err != nil {
			return err
		}
	}

	if err := c.FillTrieWithWords(ctx, ""); err != nil {
		return err
	}

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-exit
		h.Shutdown()
	}()
	return h.Listen(fmt.Sprintf(":%s", app.cfg.App.Port))
}
