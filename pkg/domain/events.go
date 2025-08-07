package domain

type (
	PaymentCallbackEvent  struct{}
	MutationCallbackEvent struct{}
	CueEvent              struct {
		Cue string `json:"cue"`
	}
)
