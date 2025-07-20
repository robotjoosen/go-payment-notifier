package sound

import "net"

type Options struct {
	ip          net.IP
	port        int
	paymentCue  string
	mutationCue string
}

type OptionFunc func(options *Options)

func getDefaultOptions() Options {
	return Options{
		ip:          net.ParseIP("127.0.0.1"),
		port:        8765,
		paymentCue:  "1",
		mutationCue: "2",
	}
}

func WithIP(ip string) OptionFunc {
	return func(options *Options) {
		options.ip = net.ParseIP(ip)
	}
}

func WithPort(port int) OptionFunc {
	return func(options *Options) {
		options.port = port
	}
}

func WithPaymentCue(cue string) OptionFunc {
	return func(options *Options) {
		options.paymentCue = cue
	}
}

func WithMutationCue(cue string) OptionFunc {
	return func(options *Options) {
		options.mutationCue = cue
	}
}
