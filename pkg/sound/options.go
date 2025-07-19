package sound

import "net"

type Options struct {
	ip      net.IP
	port    int
	address string
}

type OptionFunc func(options *Options)

func getDefaultOptions() Options {
	return Options{
		ip:      net.ParseIP("127.0.0.1"),
		port:    8765,
		address: "/osc/address",
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

func WithAddress(path string) OptionFunc {
	return func(options *Options) {
		options.address = path
	}
}
