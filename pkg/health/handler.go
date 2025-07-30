package health

import (
	"net/http"

	"github.com/robotjoosen/go-payment-notifier/pkg/server"
)

type Controller struct{}

type Response struct {
	Status string `json:"status"`
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.SuccessResponse(w, Response{Status: "OK"})
		default:
			server.CodeResponse(w, http.StatusMethodNotAllowed)
		}
	}
}
