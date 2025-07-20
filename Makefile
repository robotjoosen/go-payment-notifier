.PHONY: binary-build
binary-build:
	go build -o ./bin/notifier ./cmd/notifier && \
	go build -o ./bin/config ./cmd/config

.PHONY: docker-build
docker-build:
	 docker build --build-arg BUILD_TARGET=./cmd/notifier -t payment-notifier/notifier --platform linux/amd64 .
	 docker build --build-arg BUILD_TARGET=./cmd/config -t payment-notifier/config --platform linux/amd64 .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 --env-file .env payment-notifier/notifier 
