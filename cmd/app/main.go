package main

import (
	"context"
	"flag"
	"fmt"
	"leetgo/config"
	"leetgo/internal/app"
	"leetgo/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	log := slog.New(logger.CustomFormat(os.Stdout))

	seed := flag.Bool("seed", false, "Seed the database from TXT files in ./data/")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		log.Error(fmt.Sprintf("Config reading: %s", err))
		os.Exit(1)
	}

	a := app.New(cfg, log, *seed)
	if err = a.Run(ctx); err != nil {
		log.Error(fmt.Sprintf("App run: %s", err))
	}
}
