package shutdown

import (
	"context"
	"net/http"

	"github.com/robotjoosen/go-payment-notifier/pkg/server"
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
