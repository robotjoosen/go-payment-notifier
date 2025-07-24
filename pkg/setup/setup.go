package setup

import (
	"log/slog"
	"os"

	"github.com/OGKevin/go-bunq/bunq"
	"gitlab.com/sir-this-is-a-wendys/go-payment-notifier/pkg/env"
)

const (
	modeDevelopment = "DEV"

	defaultMode           = modeDevelopment
	defaultLogLevel       = "INFO"
	defaultServerPort     = 8080
	defaultServiceName    = "payment-notifier"
	defaultOSCAddress     = "/osc/address"
	defaultOSCIP          = "127.0.0.1"
	defaultOSCPort        = 8765
	defaultOSCPaymentCue  = 1
	defaultOSCMutationCue = 2
	defaultBunqBaseURL    = bunq.BaseURLProduction
	defaultBunqIPRange    = "185.40.108.0/22"
)

type Environment struct {
	Mode           string     `mapstructure:"MODE"`
	LogLevel       slog.Level `mapstructure:"LOG_LEVEL"`
	ServiceName    string     `mapstructure:"SERVICE_NAME"`
	ServerPort     int        `mapstructure:"SERVER_PORT"`
	BunqBaseURL    string     `mapstructure:"BUNQ_BASE_URL"`
	BunqAPIKey     string     `mapstructure:"BUNQ_API_KEY"`
	BunqAppName    string     `mapstructure:"BUNQ_APP_NAME"`
	BunqIPRange    string     `mapstructure:"BUNQ_IP_RANGE"`
	OscIP          string     `mapstructure:"OSC_IP"`
	OscPort        int        `mapstructure:"OSC_PORT"`
	OscPaymentCue  string     `mapstructure:"OSC_PAYMENT_CUE"`
	OscMutationCue string     `mapstructure:"OSC_MUTATION_CUE"`
}

func InitLog(level slog.Level) {
	hostname, err := os.Hostname()
	if err != nil {
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).
		With(
			slog.String("hostname", hostname),
		))
}

func LoadEnv() Environment {
	e, err := env.Load[Environment](map[string]any{
		"MODE":             defaultMode,
		"LOG_LEVEL":        defaultLogLevel,
		"SERVICE_NAME":     defaultServiceName,
		"SERVER_PORT":      defaultServerPort,
		"BUNQ_BASE_URL":    defaultBunqBaseURL,
		"BUNQ_API_KEY":     "",
		"BUNQ_APP_NAME":    defaultServiceName,
		"BUNQ_IP_RANGE":    defaultBunqIPRange,
		"OSC_IP":           defaultOSCIP,
		"OSC_PORT":         defaultOSCPort,
		"OSC_PAYMENT_CUE":  defaultOSCPaymentCue,
		"OSC_MUTATION_CUE": defaultOSCMutationCue,
	},
		func(e *Environment) {
			if e.Mode == modeDevelopment && e.BunqBaseURL == "" {
				e.BunqBaseURL = bunq.BaseURLProduction
			}
		},
	)
	if err != nil {
		slog.Error("failed to load environment", "err", err.Error())

		os.Exit(1)
	}

	return e
}
