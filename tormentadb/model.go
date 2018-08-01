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
	GetCreated() time.Time
	SetID(gouuidv6.UUID)
}

type Model struct {
	ID          gouuidv6.UUID `msg:",extension" tormenta:"noindex"`
	Created     time.Time     `msg:"-" tormenta:"noindex"`
	LastUpdated time.Time     `tormenta:"noindex"`
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

func (m *Model) SetID(id gouuidv6.UUID) {
	m.ID = id
}

func (m *Model) GetCreated() time.Time {
	createdAt := m.ID.Time()
	m.Created = createdAt
	return createdAt
}
