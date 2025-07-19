package broadcaster

import (
	"context"

	"github.com/anthdm/hollywood/actor"
)

type Actor struct {
	adapter IOAware
}

// NewActor creates a new actor which can be added to the broadcast engine
func NewActor(adapter IOAware) *Actor {
	return &Actor{adapter: adapter}
}

func (a *Actor) receiverFunc() func() actor.Receiver {
	return func() actor.Receiver {
		return a
	}
}

// Run listen for message coming in a broadcast events
func (a *Actor) Run(ctx context.Context, engine *actor.Engine) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-a.adapter.Output():
			engine.BroadcastEvent(msg)
		}
	}
}

// Receive magic function for handling incoming messages send to the engine
func (a *Actor) Receive(ctx *actor.Context) {
	a.adapter.Input(ctx.Message())
}
