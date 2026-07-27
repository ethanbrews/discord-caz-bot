package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethanbrews/discord-caz-bot/bot"
	"github.com/ethanbrews/discord-caz-bot/configuration"
	"github.com/monzo/terrors"
	"golang.org/x/sync/errgroup"
)

func PrintTerror(err error) {
	var terror *terrors.Error
	if errors.As(terrors.Propagate(err), &terror); terror == nil {
		return
	}
	fmt.Println(terror.ErrorMessage())
	for k, v := range terror.Params {
		fmt.Printf("\t%s=%s", k, v)
	}
	fmt.Println()
}

func main() {
	err := run()
	if err != nil {
		slog.Error("Crashed with error")
		PrintTerror(err)
	}
	syscall.Exit(1)
}

func run() error {
	configFileLocation, ok := os.LookupEnv("DISCORD_CAZ_CONFIG")
	if !ok {
		configFileLocation = "./config.json"
		slog.Warn("DISCORD_CAZ_CONFIG environment variable not set", "default", configFileLocation)
	}

	applicationConfig, err := configuration.ReadApplicationConfig(configFileLocation)
	if err != nil {
		return terrors.Augment(err, "loading config", nil)
	}

	session, err := bot.MakeBot(applicationConfig.AuthenticationToken)
	if err != nil {
		return terrors.Augment(err, "making bot", nil)
	}

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(sigCtx)

	g.Go(func() error {
		return bot.StartBot(ctx, session, applicationConfig.GuildID)
	})

	botErr := g.Wait()

	// Try to close discord bot gracefully with a 5 sec limit
	botShutdownCtx, botShutdown := context.WithCancel(context.Background())

	go func() {
		session.Close()
		botShutdown()
	}()

	select {
	case <-botShutdownCtx.Done():
	case <-time.After(5 * time.Second):
	}

	return botErr
}
