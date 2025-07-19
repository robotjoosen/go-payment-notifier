package sound

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/hypebeast/go-osc/osc"
	"gitlab.com/sir-this-is-a-wendys/go-payment-notifier/pkg/domain"
)

type Controller struct {
	identifier    uuid.UUID
	client        *osc.Client
	address       string
	outputChannel chan any
}

func New(optionFuncs ...OptionFunc) *Controller {
	options := getDefaultOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}

	c := osc.NewClient(options.ip.String(), options.port)
	if c == nil {
		slog.Error("failed to create osc client")

		return nil
	}

	return &Controller{
		client:        c,
		address:       options.address,
		outputChannel: make(chan any, 100),
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
	case domain.PaymentCallbackEvent:
		msg := osc.NewMessage(c.address)
		msg.Append(1)

		c.client.Send(msg)
	case domain.MutationCallbackEvent:
		msg := osc.NewMessage(c.address)
		msg.Append(2)

		c.client.Send(msg)
	}
}

// Output for internal message bus communication
func (c *Controller) Output() chan any {
	return c.outputChannel
}
