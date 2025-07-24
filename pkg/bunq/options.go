package bunq

import (
	"log/slog"
	"net/netip"

	bunqclient "github.com/OGKevin/go-bunq/bunq"
)

type Options struct {
	baseURL string
	appName string
	apiKey  string
	network netip.Prefix
}

type OptionFunc func(options *Options)

func getDefaultOptions() Options {
	return Options{
		baseURL: bunqclient.BaseURLSandbox,
		appName: "payment-notifier",
		apiKey:  "",
	}
}

func WithBaseURL(url string) OptionFunc {
	return func(options *Options) {
		options.baseURL = url
	}
}

func WithAPIKey(key string) OptionFunc {
	return func(options *Options) {
		options.apiKey = key
	}
}

func WithAppName(name string) OptionFunc {
	return func(options *Options) {
		options.appName = name
	}
}

func WithIPRange(network string) OptionFunc {
	return func(options *Options) {
		network, err := netip.ParsePrefix(network)
		if err != nil {
			slog.Error("invalid ip range given", "err", err.Error())

			return
		}

		options.network = network
	}
}
