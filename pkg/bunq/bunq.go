package bunq

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	bunqclient "github.com/OGKevin/go-bunq/bunq"
	"github.com/google/uuid"
	"gitlab.com/sir-this-is-a-wendys/go-payment-notifier/pkg/domain"
	"gitlab.com/sir-this-is-a-wendys/go-payment-notifier/pkg/server"
)

type Controller struct {
	identifier    uuid.UUID
	outputChannel chan any
	client        *bunqclient.Client
}

const (
	callbackPathPayment  = "payment/callback"
	callbackPathMutation = "mutation/callback"
)

func New(optionFuncs ...OptionFunc) *Controller {
	options := getDefaultOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}

	// connect tp bunq
	key, err := bunqclient.CreateNewKeyPair()
	if err != nil {
		slog.Error("failed to new key pair", "err", err.Error())

		return nil
	}

	c := bunqclient.NewClient(
		context.Background(),
		options.baseURL,
		key,
		options.apiKey,
		options.appName,
	)
	if err := c.Init(); err != nil {
		slog.Error("failed to initialize bunq client", "err", err.Error())

		return nil
	}

	return &Controller{
		outputChannel: make(chan any, 100),
		client:        c,
	}
}

func (c *Controller) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, callbackPathPayment):
			c.outputChannel <- domain.PaymentCallbackEvent{} // TODO: add payment data for additional event handling
		case strings.Contains(r.URL.Path, callbackPathMutation):
			c.outputChannel <- domain.MutationCallbackEvent{} // TODO: add mutaiton data for addiontal event handling
		}

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
