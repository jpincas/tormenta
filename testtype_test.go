package tormenta_test

import (
	"errors"

	"github.com/jpincas/tormenta"
)

type (
	DefinedInt    int
	DefinedString string
	DefinedFloat  float64
	DefinedBool   bool
)

type EmbeddedStruct struct {
	EmbeddedIntField    int
	EmbeddedStringField string
	EmbeddedFloatField  float64
	EmbeddedBoolField   bool
}

type TestType struct {
	tormenta.Model

	// Basic types
	IntField          int
	StringField       string
	MultipleWordField string `tormenta:"split"`
	FloatField        float64
	BoolField         bool

	// Slice types
	IntSliceField    []int
	StringSliceField []string
	FloatSliceField  []float64
	BoolSliceField   []bool

	// Defined types
	DefinedIntField    DefinedInt
	DefinedStringField DefinedString
	DefinedFloatField  DefinedFloat
	DefinedBoolField   DefinedBool

	// Defined slice types
	DefinedIntSliceField    []DefinedInt
	DefinedStringSliceField []DefinedString
	DefinedFloatSliceField  []DefinedFloat
	DefinedBoolSliceField   []DefinedBool

	// Embedded struct
	EmbeddedStruct

	// Fields for trigger testing
	TriggerString   string
	Retrieved       bool
	IsSaved         bool
	ShouldBlockSave bool
}

func (t TestType) PreSave() error {
	if t.ShouldBlockSave {
		return errors.New("presave trigger is blocking save")
	}

	return nil
}

func (t *TestType) PostSave() {
	t.IsSaved = true
}

func (t *TestType) PostGet(ctx map[string]interface{}) {
	sessionIdFromContext, ok := ctx["sessionid"]
	if ok {
		if sessionId, ok := sessionIdFromContext.(string); ok {
			t.TriggerString = sessionId
		}
	}

	t.Retrieved = true
}

type TestType2 struct {
	tormenta.Model

	IntField    int
	StringField string
	FloatField  float64
	BoolField   bool
}
