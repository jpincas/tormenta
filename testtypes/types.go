package testtypes

import (
	"errors"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

type (
	DefinedID     gouuidv6.UUID
	DefinedInt    int
	DefinedString string
	DefinedFloat  float64
	DefinedBool   bool
	DefinedDate   time.Time
)

type MyStruct struct {
	StructIntField    int
	StructStringField string
	StructFloatField  float64
	StructBoolField   bool
	StructDateField   time.Time
}

type FullStruct struct {
	tormenta.Model

	// Basic types
	IntField          int
	StringField       string
	MultipleWordField string `tormenta:"split"`
	FloatField        float64
	BoolField         bool
	DateField         time.Time

	// Slice types
	IDSliceField     []gouuidv6.UUID
	IntSliceField    []int
	StringSliceField []string
	FloatSliceField  []float64
	BoolSliceField   []bool
	DateSliceField   []time.Time

	// Defined types
	DefinedIDField     DefinedID
	DefinedIntField    DefinedInt
	DefinedStringField DefinedString
	DefinedFloatField  DefinedFloat
	DefinedBoolField   DefinedBool
	DefinedDateField   DefinedDate

	// Defined slice types
	DefinedIDSliceField     []DefinedID
	DefinedIntSliceField    []DefinedInt
	DefinedStringSliceField []DefinedString
	DefinedFloatSliceField  []DefinedFloat
	DefinedBoolSliceField   []DefinedBool
	DefinedDateSliceField   []DefinedDate

	// Embedded struct
	MyStruct

	// Named Struct
	StructField MyStruct

	// Fields for trigger testing
	TriggerString   string
	Retrieved       bool
	IsSaved         bool
	ShouldBlockSave bool
}

func (t FullStruct) PreSave() error {
	if t.ShouldBlockSave {
		return errors.New("presave trigger is blocking save")
	}

	return nil
}

func (t *FullStruct) PostSave() {
	t.IsSaved = true
}

func (t *FullStruct) PostGet(ctx map[string]interface{}) {
	sessionIdFromContext, ok := ctx["sessionid"]
	if ok {
		if sessionId, ok := sessionIdFromContext.(string); ok {
			t.TriggerString = sessionId
		}
	}

	t.Retrieved = true
}

type MiniStruct struct {
	tormenta.Model

	IntField    int
	StringField string
	FloatField  float64
	BoolField   bool
}
