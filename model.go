package tormenta

import (
	"time"

	"github.com/jpincas/gouuidv6"
)

//go:generate msgp

type Model struct {
	ID          gouuidv6.UUID `msg:",extension"`
	LastUpdated time.Time
}

func newID() gouuidv6.UUID {
	return gouuidv6.New()
}

func newModel() Model {
	return Model{
		ID:          gouuidv6.New(),
		LastUpdated: time.Now(),
	}
}
