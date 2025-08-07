package cue

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/robotjoosen/go-payment-notifier/pkg/domain"
	"github.com/robotjoosen/go-payment-notifier/pkg/server"
)

type Controller struct {
	identifier    uuid.UUID
	outputChannel chan any
	cueMap        *xsync.Map[string, string]
}

func New() *Controller {
	return &Controller{
		outputChannel: make(chan any, 100),
		cueMap:        xsync.NewMap[string, string](),
	}
}

func (c *Controller) AddEndpoint(path string, cuePath string) {
	c.cueMap.Store(path, cuePath)
}

func (c *Controller) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("got request", "url", r.URL.Path)

		cn, exists := c.cueMap.Load(r.URL.Path)
		if !exists {
			server.CodeResponse(w, http.StatusNotFound)

			return
		}

		c.outputChannel <- domain.CueEvent{
			Cue: cn,
		}

		slog.Info("dispatched event")

		server.CodeResponse(w, http.StatusAccepted)
	}
}

func (c *Controller) Identifier() uuid.UUID {
	if c.identifier == uuid.Nil {
		c.identifier = uuid.New()
	}

	return c.identifier
}

func (c *Controller) Input(msg any) {
	switch msg.(type) {
	}
}

// Output for internal message bus communication
func (c *Controller) Output() chan any {
	return c.outputChannel
}
