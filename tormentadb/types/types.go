package types

import tormenta "github.com/jpincas/tormenta/tormentadb"

//go:generate msgp

type (
	DefinedInt    int
	DefinedString string
	DefinedFloat  float64
	DefinedBool   bool
)

type TestType struct {
	tormenta.Model

	// Basic types
	IntField    int
	StringField string
	FloatField  float64
	BoolField   bool

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
}
