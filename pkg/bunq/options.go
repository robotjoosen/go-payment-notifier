package bunq

import bunqclient "github.com/OGKevin/go-bunq/bunq"

type Options struct {
	baseURL string
	appName string
	apiKey  string
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
