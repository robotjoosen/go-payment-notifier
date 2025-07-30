package main

import (
	"log/slog"
	"os"

	"github.com/robotjoosen/go-payment-notifier/pkg/bunq"
	"github.com/robotjoosen/go-payment-notifier/pkg/setup"
)

func main() {
	e := setup.LoadEnv()
	setup.InitLog(e.LogLevel)

	slog.Info("using environment", "configuration", e)

	bunqInstance := bunq.New()
	if bunqInstance == nil {
		slog.Error("failed to initialise bunq instance")

		os.Exit(1)
	}

	if err := bunqInstance.Connect(
		bunq.WithBaseURL(e.BunqBaseURL),
		bunq.WithAPIKey(e.BunqAPIKey),
		bunq.WithAppName(e.BunqAppName),
	); err != nil {
		slog.Error("failed to connect to bunq", "err", err.Error())

		os.Exit(1)
	}

	if err := bunqInstance.SetNotificationWebhook(); err != nil {
		slog.Error("failed to set notification filters", "err", err.Error())

		os.Exit(1)
	}
}
