package tormentadb

import (
	"time"

	"github.com/jpincas/gouuidv6"
)

//go:generate msgp

type Tormentable interface {
	MarshalMsg([]byte) ([]byte, error)
	UnmarshalMsg([]byte) ([]byte, error)
	PreSave() error
	PostSave()
	PostGet()
}

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

func (m Model) PreSave() error {
	return nil
}

func (m Model) PostSave() {}

func (m Model) PostGet() {}
