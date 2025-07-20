package bunq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	CallbackPathPayment  = "/callbacks/payment"
	CallbackPathMutation = "/callbacks/mutation"
)

func New() *Controller {
	// connect tp bunq

	return &Controller{
		outputChannel: make(chan any, 100),
	}
}

func (c *Controller) Connect(optionFuncs ...OptionFunc) error {
	options := getDefaultOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}

	key, err := bunqclient.CreateNewKeyPair()
	if err != nil {
		slog.Error("failed to new key pair", "err", err.Error())

		return err
	}
	c.client = bunqclient.NewClient(
		context.Background(),
		options.baseURL,
		key,
		options.apiKey,
		options.appName,
	)
	if err := c.client.Init(); err != nil {
		slog.Error("failed to initialize bunq client", "err", err.Error())

		return err
	}

	return nil
}

func (c *Controller) SetNotificationWebhook() error {
	accounts, err := c.client.AccountService.GetMonetaryAccountBank(1)
	if err != nil {
		return err
	}

	userID, err := c.client.GetUserID()
	if err != nil {
		return err
	}

	n := NotificationFilters{
		NotificationFilter: []NotificationFilter{
			{
				Category:           "PAYMENT",
				NotificationTarget: "http://127.0.0.1/payment/callback",
			},
			{
				Category:           "MUTATION",
				NotificationTarget: "http://127.0.0.1/mutation/callback",
			},
		},
	}

	rawNotificationFilters, err := json.Marshal(n)
	if err != nil {
		return err
	}

	for _, account := range accounts.Response {
		if _, err := c.client.Post(
			fmt.Sprintf(
				"/v1/user/%d/monetary-account/%d/notification-filter-url",
				userID,
				account.MonetaryAccountBank.ID,
			),
			"application/json",
			bytes.NewBuffer(rawNotificationFilters),
		); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, CallbackPathPayment):
			c.outputChannel <- domain.PaymentCallbackEvent{} // TODO: add payment data for additional event handling
		case strings.Contains(r.URL.Path, CallbackPathMutation):
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
