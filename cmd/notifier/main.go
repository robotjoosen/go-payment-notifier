package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/robotjoosen/go-payment-notifier/pkg/broadcaster"
	"github.com/robotjoosen/go-payment-notifier/pkg/bunq"
	"github.com/robotjoosen/go-payment-notifier/pkg/health"
	"github.com/robotjoosen/go-payment-notifier/pkg/server"
	"github.com/robotjoosen/go-payment-notifier/pkg/setup"
	"github.com/robotjoosen/go-payment-notifier/pkg/shutdown"
	"github.com/robotjoosen/go-payment-notifier/pkg/sound"
)

func main() {
	e := setup.LoadEnv()
	setup.InitLog(e.LogLevel)

	slog.Info("using environment", "configuration", e)

	bunqInstance := bunq.New(bunq.WithIPRange(e.BunqIPRange))
	if bunqInstance == nil {
		slog.Error("failed to initialise bunq instance")

		os.Exit(1)
	}

	soundInstance := sound.New(
		sound.WithIP(e.OscIP),
		sound.WithPort(e.OscPort),
		sound.WithPaymentCue(e.OscPaymentCue),
		sound.WithMutationCue(e.OscMutationCue),
	)
	if soundInstance == nil {
		slog.Error("failed to initialise osc instance")

		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	serverInstance := server.New(e.ServerPort, map[string]http.HandlerFunc{
		"GET /health":   health.New().Handler(),
		"GET /shutdown": shutdown.New().Handler(cancel),

		// TODO: unsure how bunq webhook are formatted,
		// 			the most generic endpoint is used in this case.
		// 			improvements might be required
		bunq.CallbackPathPayment:  bunqInstance.Handler(),
		bunq.CallbackPathMutation: bunqInstance.Handler(),
	})
	go serverInstance.Run()

	// connect services to internal message bus
	actorEngine := broadcaster.New()
	if actorEngine == nil {
		slog.Error("failed to create broadcaster")

		os.Exit(1)
	}
	actorEngine.BulkAdd(
		bunqInstance,
		soundInstance,
	)

	// shutdown on context cancellation
	shutdownWait(ctx, []func(){
		serverInstance.Stop,
	})
}

func shutdownWait(parentCtx context.Context, shutdownFuncs []func()) {
	ctx, ctxStop := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM)
	shutdownFuncs = append(shutdownFuncs, ctxStop)

	<-ctx.Done()

	slog.Info("shutdown started", "process_total", len(shutdownFuncs))
	for _, shutdownFunc := range shutdownFuncs {
		shutdownFunc()
	}

	slog.Info("shutdown finished")
}
