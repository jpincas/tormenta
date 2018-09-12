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
	PostGet(ctx map[string]interface{})
	GetCreated() time.Time
	SetID(gouuidv6.UUID)
}

type Model struct {
	ID          gouuidv6.UUID `msg:",extension" tormenta:"noindex" json:"id"`
	Created     time.Time     `msg:"-" tormenta:"noindex" json:"created"`
	LastUpdated time.Time     `json:"lastUpdated"`
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

func (m Model) PostGet(ctx map[string]interface{}) {}

func (m *Model) SetID(id gouuidv6.UUID) {
	m.ID = id
}

func (m *Model) GetCreated() time.Time {
	createdAt := m.ID.Time()
	m.Created = createdAt
	return createdAt
}
