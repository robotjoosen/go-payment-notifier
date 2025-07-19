.PHONY: build
build:
	go build -o ./bin/notifier ./cmd/notifier && \
	go build -o ./bin/config ./cmd/config

