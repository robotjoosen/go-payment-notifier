package broadcaster

import (
	"context"
	"sync"

	"github.com/anthdm/hollywood/actor"
	"github.com/google/uuid"
)

type IOAware interface {
	Input(interface{})
	Output() chan interface{}
	Identifier() uuid.UUID
}

type Service struct {
	actors []*Actor
	mux    sync.Mutex
	engine *actor.Engine
}

// New internal message bus
func New() *Service {
	e, err := actor.NewEngine(actor.EngineConfig{})
	if err != nil {
		return nil
	}

	return &Service{
		engine: e,
		mux:    sync.Mutex{},
		actors: make([]*Actor, 0),
	}
}

// BulkAdd multiple In/Output aware services to the broadcast engine
func (s *Service) BulkAdd(ios ...IOAware) *Service {
	for _, io := range ios {
		s.Add(io)
	}

	return s
}

// Add single In/Output aware service to the broadcast engine
func (s *Service) Add(io IOAware) *Service {
	s.mux.Lock()
	defer s.mux.Unlock()

	act := NewActor(io)
	s.actors = append(s.actors, act)

	go act.Run(context.Background(), s.engine) // todo: figure out if we need context cancelling

	if a := s.engine.Spawn(
		NewActor(io).receiverFunc(),
		"default",
		actor.WithID(io.Identifier().String()),
	); a != nil {
		s.engine.Subscribe(a)
	}

	return s
}
