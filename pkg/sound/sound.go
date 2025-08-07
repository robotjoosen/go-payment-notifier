package sound

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hypebeast/go-osc/osc"
	"github.com/robotjoosen/go-payment-notifier/pkg/domain"
)

const cueStartPath = "/cue/%s/start"

type Controller struct {
	identifier    uuid.UUID
	client        *osc.Client
	paymentCue    string
	mutationCue   string
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
		paymentCue:    options.paymentCue,
		mutationCue:   options.mutationCue,
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
	switch m := msg.(type) {
	case domain.CueEvent:
		msg := osc.NewMessage(m.Cue)

		if err := c.client.Send(msg); err != nil {
			slog.Error("send failed",
				slog.String("err", err.Error()),
				slog.Any("msg", msg),
			)

			return
		}

		slog.Info("Event triggered",
			slog.Any("msg", msg),
		)

	case domain.PaymentCallbackEvent:
		msg := osc.NewMessage(fmt.Sprintf(cueStartPath, c.paymentCue))

		if err := c.client.Send(msg); err != nil {
			slog.Error("send failed",
				slog.String("err", err.Error()),
				slog.Any("msg", msg),
			)

			return
		}

		slog.Info("Payment callback event triggered",
			slog.Any("msg", msg),
		)
	case domain.MutationCallbackEvent:
		msg := osc.NewMessage(fmt.Sprintf(cueStartPath, c.mutationCue))

		if err := c.client.Send(msg); err != nil {
			slog.Error("send failed",
				slog.String("err", err.Error()),
				slog.Any("msg", msg),
			)

			return
		}

		slog.Info("Mutation callback triggered",
			slog.Any("msg", msg),
		)
	}
}

// Output for internal message bus communication
func (c *Controller) Output() chan any {
	return c.outputChannel
}
