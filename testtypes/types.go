package testtypes

import (
	"errors"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

//go:generate ffjson $GOFILE

type (
	DefinedID     gouuidv6.UUID
	DefinedInt    int
	DefinedInt16  int16
	DefinedUint16 uint16
	DefinedString string
	DefinedFloat  float64
	DefinedBool   bool
	DefinedDate   time.Time
	DefinedStruct MyStruct
)

type MyStruct struct {
	StructIntField    int
	StructStringField string
	StructFloatField  float64
	StructBoolField   bool
	StructDateField   time.Time
}

type RelatedStruct struct {
	tormenta.Model

	StructIntField    int
	StructStringField string
	StructFloatField  float64
	StructBoolField   bool
	StructDateField   time.Time

	NestedID gouuidv6.UUID
	Nested   *NestedRelatedStruct `tormenta:"-"`

	// For 'belongs to'
	FullStructID gouuidv6.UUID
}

type NestedRelatedStruct struct {
	tormenta.Model

	NestedID gouuidv6.UUID
	Nested   *DoubleNestedRelatedStruct `tormenta:"-"`
}

type DoubleNestedRelatedStruct struct {
	tormenta.Model
}

type FullStruct struct {
	tormenta.Model

	// Basic types
	IntField          int
	IDField           gouuidv6.UUID
	AnotherIntField   int
	StringField       string
	MultipleWordField string `tormenta:"split"`
	FloatField        float64
	Float32Field      float32
	AnotherFloatField float64
	BoolField         bool
	DateField         time.Time

	// Fixed-length types
	Int8Field  int8
	Int16Field int16
	Int32Field int32
	Int64Field int64

	UintField   uint
	Uint8Field  uint8
	Uint16Field uint16
	Uint32Field uint32
	Uint64Field uint64

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
	DefinedStructField DefinedStruct

	// Defined Fixed-length types - just a sample
	DefinedInt16Field  DefinedInt16
	DefinedUint16Field DefinedUint16

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
	StructField MyStruct `tormenta:"nested"`

	// Fields for trigger testing
	TriggerString   string
	Retrieved       bool
	IsSaved         bool
	ShouldBlockSave bool

	// Fields for 'no index' testing
	NoIndexSimple                string `tormenta:"noindex"`
	NoIndexTwoTags               string `tormenta:"noindex; split"`
	NoIndexTwoTagsDifferentOrder string `tormenta:"split;noindex"`

	// Fields for 'no save' testing
	NoSaveSimple                string `tormenta:"-"`
	NoSaveTwoTags               string `tormenta:"split;-"`
	NoSaveTwoTagsDifferentOrder string `tormenta:"-;split"`

	// for this one we change the field name with a json tag
	NoSaveJSONtag     string `tormenta:"-" json:"noSaveJsonTag"`
	NoSaveJSONSkiptag string `tormenta:"-" json:"-"`

	// Fields for relations testing
	HasOneID gouuidv6.UUID
	HasOne   *RelatedStruct `tormenta:"-"`

	HasAnotherOneID gouuidv6.UUID
	HasAnotherOne   *RelatedStruct `tormenta:"-"`

	HasManyIDs []gouuidv6.UUID
	HasMany    []*RelatedStruct `tormenta:"-"`

	RelatedStructsByQuery []*RelatedStruct `tormenta:"-"`
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
