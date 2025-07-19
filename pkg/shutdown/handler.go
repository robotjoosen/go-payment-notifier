package shutdown

import (
	"context"
	"net/http"

	"gitlab.com/sir-this-is-a-wendys/go-payment-notifier/pkg/server"
)

type Handler struct{}

func New() Handler {
	return Handler{}
}

func (h Handler) Handler(c context.CancelFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		c()

		server.CodeResponse(w, http.StatusAccepted)
	}
}
