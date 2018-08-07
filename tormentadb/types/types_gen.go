package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *DefinedBool) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 bool
		zb0001, err = dc.ReadBool()
		if err != nil {
			return
		}
		(*z) = DefinedBool(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DefinedBool) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBool(bool(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DefinedBool) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBool(o, bool(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DefinedBool) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 bool
		zb0001, bts, err = msgp.ReadBoolBytes(bts)
		if err != nil {
			return
		}
		(*z) = DefinedBool(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z DefinedBool) Msgsize() (s int) {
	s = msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DefinedFloat) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 float64
		zb0001, err = dc.ReadFloat64()
		if err != nil {
			return
		}
		(*z) = DefinedFloat(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DefinedFloat) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteFloat64(float64(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DefinedFloat) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendFloat64(o, float64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DefinedFloat) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 float64
		zb0001, bts, err = msgp.ReadFloat64Bytes(bts)
		if err != nil {
			return
		}
		(*z) = DefinedFloat(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z DefinedFloat) Msgsize() (s int) {
	s = msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DefinedInt) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 int
		zb0001, err = dc.ReadInt()
		if err != nil {
			return
		}
		(*z) = DefinedInt(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DefinedInt) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteInt(int(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DefinedInt) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt(o, int(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DefinedInt) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 int
		zb0001, bts, err = msgp.ReadIntBytes(bts)
		if err != nil {
			return
		}
		(*z) = DefinedInt(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z DefinedInt) Msgsize() (s int) {
	s = msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DefinedString) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 string
		zb0001, err = dc.ReadString()
		if err != nil {
			return
		}
		(*z) = DefinedString(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z DefinedString) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteString(string(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z DefinedString) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendString(o, string(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DefinedString) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 string
		zb0001, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
		(*z) = DefinedString(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z DefinedString) Msgsize() (s int) {
	s = msgp.StringPrefixSize + len(string(z))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TestType) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Model":
			err = z.Model.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "IntField":
			z.IntField, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "StringField":
			z.StringField, err = dc.ReadString()
			if err != nil {
				return
			}
		case "FloatField":
			z.FloatField, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "BoolField":
			z.BoolField, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "IntSliceField":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.IntSliceField) >= int(zb0002) {
				z.IntSliceField = (z.IntSliceField)[:zb0002]
			} else {
				z.IntSliceField = make([]int, zb0002)
			}
			for za0001 := range z.IntSliceField {
				z.IntSliceField[za0001], err = dc.ReadInt()
				if err != nil {
					return
				}
			}
		case "StringSliceField":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.StringSliceField) >= int(zb0003) {
				z.StringSliceField = (z.StringSliceField)[:zb0003]
			} else {
				z.StringSliceField = make([]string, zb0003)
			}
			for za0002 := range z.StringSliceField {
				z.StringSliceField[za0002], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "FloatSliceField":
			var zb0004 uint32
			zb0004, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.FloatSliceField) >= int(zb0004) {
				z.FloatSliceField = (z.FloatSliceField)[:zb0004]
			} else {
				z.FloatSliceField = make([]float64, zb0004)
			}
			for za0003 := range z.FloatSliceField {
				z.FloatSliceField[za0003], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "BoolSliceField":
			var zb0005 uint32
			zb0005, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.BoolSliceField) >= int(zb0005) {
				z.BoolSliceField = (z.BoolSliceField)[:zb0005]
			} else {
				z.BoolSliceField = make([]bool, zb0005)
			}
			for za0004 := range z.BoolSliceField {
				z.BoolSliceField[za0004], err = dc.ReadBool()
				if err != nil {
					return
				}
			}
		case "DefinedIntField":
			{
				var zb0006 int
				zb0006, err = dc.ReadInt()
				if err != nil {
					return
				}
				z.DefinedIntField = DefinedInt(zb0006)
			}
		case "DefinedStringField":
			{
				var zb0007 string
				zb0007, err = dc.ReadString()
				if err != nil {
					return
				}
				z.DefinedStringField = DefinedString(zb0007)
			}
		case "DefinedFloatField":
			{
				var zb0008 float64
				zb0008, err = dc.ReadFloat64()
				if err != nil {
					return
				}
				z.DefinedFloatField = DefinedFloat(zb0008)
			}
		case "DefinedBoolField":
			{
				var zb0009 bool
				zb0009, err = dc.ReadBool()
				if err != nil {
					return
				}
				z.DefinedBoolField = DefinedBool(zb0009)
			}
		case "DefinedIntSliceField":
			var zb0010 uint32
			zb0010, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DefinedIntSliceField) >= int(zb0010) {
				z.DefinedIntSliceField = (z.DefinedIntSliceField)[:zb0010]
			} else {
				z.DefinedIntSliceField = make([]DefinedInt, zb0010)
			}
			for za0005 := range z.DefinedIntSliceField {
				{
					var zb0011 int
					zb0011, err = dc.ReadInt()
					if err != nil {
						return
					}
					z.DefinedIntSliceField[za0005] = DefinedInt(zb0011)
				}
			}
		case "DefinedStringSliceField":
			var zb0012 uint32
			zb0012, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DefinedStringSliceField) >= int(zb0012) {
				z.DefinedStringSliceField = (z.DefinedStringSliceField)[:zb0012]
			} else {
				z.DefinedStringSliceField = make([]DefinedString, zb0012)
			}
			for za0006 := range z.DefinedStringSliceField {
				{
					var zb0013 string
					zb0013, err = dc.ReadString()
					if err != nil {
						return
					}
					z.DefinedStringSliceField[za0006] = DefinedString(zb0013)
				}
			}
		case "DefinedFloatSliceField":
			var zb0014 uint32
			zb0014, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DefinedFloatSliceField) >= int(zb0014) {
				z.DefinedFloatSliceField = (z.DefinedFloatSliceField)[:zb0014]
			} else {
				z.DefinedFloatSliceField = make([]DefinedFloat, zb0014)
			}
			for za0007 := range z.DefinedFloatSliceField {
				{
					var zb0015 float64
					zb0015, err = dc.ReadFloat64()
					if err != nil {
						return
					}
					z.DefinedFloatSliceField[za0007] = DefinedFloat(zb0015)
				}
			}
		case "DefinedBoolSliceField":
			var zb0016 uint32
			zb0016, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DefinedBoolSliceField) >= int(zb0016) {
				z.DefinedBoolSliceField = (z.DefinedBoolSliceField)[:zb0016]
			} else {
				z.DefinedBoolSliceField = make([]DefinedBool, zb0016)
			}
			for za0008 := range z.DefinedBoolSliceField {
				{
					var zb0017 bool
					zb0017, err = dc.ReadBool()
					if err != nil {
						return
					}
					z.DefinedBoolSliceField[za0008] = DefinedBool(zb0017)
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TestType) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 17
	// write "Model"
	err = en.Append(0xde, 0x0, 0x11, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	if err != nil {
		return
	}
	err = z.Model.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "IntField"
	err = en.Append(0xa8, 0x49, 0x6e, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteInt(z.IntField)
	if err != nil {
		return
	}
	// write "StringField"
	err = en.Append(0xab, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.StringField)
	if err != nil {
		return
	}
	// write "FloatField"
	err = en.Append(0xaa, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.FloatField)
	if err != nil {
		return
	}
	// write "BoolField"
	err = en.Append(0xa9, 0x42, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.BoolField)
	if err != nil {
		return
	}
	// write "IntSliceField"
	err = en.Append(0xad, 0x49, 0x6e, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.IntSliceField)))
	if err != nil {
		return
	}
	for za0001 := range z.IntSliceField {
		err = en.WriteInt(z.IntSliceField[za0001])
		if err != nil {
			return
		}
	}
	// write "StringSliceField"
	err = en.Append(0xb0, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.StringSliceField)))
	if err != nil {
		return
	}
	for za0002 := range z.StringSliceField {
		err = en.WriteString(z.StringSliceField[za0002])
		if err != nil {
			return
		}
	}
	// write "FloatSliceField"
	err = en.Append(0xaf, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.FloatSliceField)))
	if err != nil {
		return
	}
	for za0003 := range z.FloatSliceField {
		err = en.WriteFloat64(z.FloatSliceField[za0003])
		if err != nil {
			return
		}
	}
	// write "BoolSliceField"
	err = en.Append(0xae, 0x42, 0x6f, 0x6f, 0x6c, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.BoolSliceField)))
	if err != nil {
		return
	}
	for za0004 := range z.BoolSliceField {
		err = en.WriteBool(z.BoolSliceField[za0004])
		if err != nil {
			return
		}
	}
	// write "DefinedIntField"
	err = en.Append(0xaf, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteInt(int(z.DefinedIntField))
	if err != nil {
		return
	}
	// write "DefinedStringField"
	err = en.Append(0xb2, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(string(z.DefinedStringField))
	if err != nil {
		return
	}
	// write "DefinedFloatField"
	err = en.Append(0xb1, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteFloat64(float64(z.DefinedFloatField))
	if err != nil {
		return
	}
	// write "DefinedBoolField"
	err = en.Append(0xb0, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x42, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(bool(z.DefinedBoolField))
	if err != nil {
		return
	}
	// write "DefinedIntSliceField"
	err = en.Append(0xb4, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.DefinedIntSliceField)))
	if err != nil {
		return
	}
	for za0005 := range z.DefinedIntSliceField {
		err = en.WriteInt(int(z.DefinedIntSliceField[za0005]))
		if err != nil {
			return
		}
	}
	// write "DefinedStringSliceField"
	err = en.Append(0xb7, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.DefinedStringSliceField)))
	if err != nil {
		return
	}
	for za0006 := range z.DefinedStringSliceField {
		err = en.WriteString(string(z.DefinedStringSliceField[za0006]))
		if err != nil {
			return
		}
	}
	// write "DefinedFloatSliceField"
	err = en.Append(0xb6, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.DefinedFloatSliceField)))
	if err != nil {
		return
	}
	for za0007 := range z.DefinedFloatSliceField {
		err = en.WriteFloat64(float64(z.DefinedFloatSliceField[za0007]))
		if err != nil {
			return
		}
	}
	// write "DefinedBoolSliceField"
	err = en.Append(0xb5, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x42, 0x6f, 0x6f, 0x6c, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.DefinedBoolSliceField)))
	if err != nil {
		return
	}
	for za0008 := range z.DefinedBoolSliceField {
		err = en.WriteBool(bool(z.DefinedBoolSliceField[za0008]))
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TestType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 17
	// string "Model"
	o = append(o, 0xde, 0x0, 0x11, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	o, err = z.Model.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "IntField"
	o = append(o, 0xa8, 0x49, 0x6e, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendInt(o, z.IntField)
	// string "StringField"
	o = append(o, 0xab, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendString(o, z.StringField)
	// string "FloatField"
	o = append(o, 0xaa, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendFloat64(o, z.FloatField)
	// string "BoolField"
	o = append(o, 0xa9, 0x42, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendBool(o, z.BoolField)
	// string "IntSliceField"
	o = append(o, 0xad, 0x49, 0x6e, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.IntSliceField)))
	for za0001 := range z.IntSliceField {
		o = msgp.AppendInt(o, z.IntSliceField[za0001])
	}
	// string "StringSliceField"
	o = append(o, 0xb0, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.StringSliceField)))
	for za0002 := range z.StringSliceField {
		o = msgp.AppendString(o, z.StringSliceField[za0002])
	}
	// string "FloatSliceField"
	o = append(o, 0xaf, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.FloatSliceField)))
	for za0003 := range z.FloatSliceField {
		o = msgp.AppendFloat64(o, z.FloatSliceField[za0003])
	}
	// string "BoolSliceField"
	o = append(o, 0xae, 0x42, 0x6f, 0x6f, 0x6c, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.BoolSliceField)))
	for za0004 := range z.BoolSliceField {
		o = msgp.AppendBool(o, z.BoolSliceField[za0004])
	}
	// string "DefinedIntField"
	o = append(o, 0xaf, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendInt(o, int(z.DefinedIntField))
	// string "DefinedStringField"
	o = append(o, 0xb2, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendString(o, string(z.DefinedStringField))
	// string "DefinedFloatField"
	o = append(o, 0xb1, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendFloat64(o, float64(z.DefinedFloatField))
	// string "DefinedBoolField"
	o = append(o, 0xb0, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x42, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendBool(o, bool(z.DefinedBoolField))
	// string "DefinedIntSliceField"
	o = append(o, 0xb4, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DefinedIntSliceField)))
	for za0005 := range z.DefinedIntSliceField {
		o = msgp.AppendInt(o, int(z.DefinedIntSliceField[za0005]))
	}
	// string "DefinedStringSliceField"
	o = append(o, 0xb7, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DefinedStringSliceField)))
	for za0006 := range z.DefinedStringSliceField {
		o = msgp.AppendString(o, string(z.DefinedStringSliceField[za0006]))
	}
	// string "DefinedFloatSliceField"
	o = append(o, 0xb6, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DefinedFloatSliceField)))
	for za0007 := range z.DefinedFloatSliceField {
		o = msgp.AppendFloat64(o, float64(z.DefinedFloatSliceField[za0007]))
	}
	// string "DefinedBoolSliceField"
	o = append(o, 0xb5, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x42, 0x6f, 0x6f, 0x6c, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DefinedBoolSliceField)))
	for za0008 := range z.DefinedBoolSliceField {
		o = msgp.AppendBool(o, bool(z.DefinedBoolSliceField[za0008]))
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TestType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Model":
			bts, err = z.Model.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "IntField":
			z.IntField, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "StringField":
			z.StringField, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "FloatField":
			z.FloatField, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "BoolField":
			z.BoolField, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "IntSliceField":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.IntSliceField) >= int(zb0002) {
				z.IntSliceField = (z.IntSliceField)[:zb0002]
			} else {
				z.IntSliceField = make([]int, zb0002)
			}
			for za0001 := range z.IntSliceField {
				z.IntSliceField[za0001], bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					return
				}
			}
		case "StringSliceField":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.StringSliceField) >= int(zb0003) {
				z.StringSliceField = (z.StringSliceField)[:zb0003]
			} else {
				z.StringSliceField = make([]string, zb0003)
			}
			for za0002 := range z.StringSliceField {
				z.StringSliceField[za0002], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "FloatSliceField":
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.FloatSliceField) >= int(zb0004) {
				z.FloatSliceField = (z.FloatSliceField)[:zb0004]
			} else {
				z.FloatSliceField = make([]float64, zb0004)
			}
			for za0003 := range z.FloatSliceField {
				z.FloatSliceField[za0003], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "BoolSliceField":
			var zb0005 uint32
			zb0005, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.BoolSliceField) >= int(zb0005) {
				z.BoolSliceField = (z.BoolSliceField)[:zb0005]
			} else {
				z.BoolSliceField = make([]bool, zb0005)
			}
			for za0004 := range z.BoolSliceField {
				z.BoolSliceField[za0004], bts, err = msgp.ReadBoolBytes(bts)
				if err != nil {
					return
				}
			}
		case "DefinedIntField":
			{
				var zb0006 int
				zb0006, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					return
				}
				z.DefinedIntField = DefinedInt(zb0006)
			}
		case "DefinedStringField":
			{
				var zb0007 string
				zb0007, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				z.DefinedStringField = DefinedString(zb0007)
			}
		case "DefinedFloatField":
			{
				var zb0008 float64
				zb0008, bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
				z.DefinedFloatField = DefinedFloat(zb0008)
			}
		case "DefinedBoolField":
			{
				var zb0009 bool
				zb0009, bts, err = msgp.ReadBoolBytes(bts)
				if err != nil {
					return
				}
				z.DefinedBoolField = DefinedBool(zb0009)
			}
		case "DefinedIntSliceField":
			var zb0010 uint32
			zb0010, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DefinedIntSliceField) >= int(zb0010) {
				z.DefinedIntSliceField = (z.DefinedIntSliceField)[:zb0010]
			} else {
				z.DefinedIntSliceField = make([]DefinedInt, zb0010)
			}
			for za0005 := range z.DefinedIntSliceField {
				{
					var zb0011 int
					zb0011, bts, err = msgp.ReadIntBytes(bts)
					if err != nil {
						return
					}
					z.DefinedIntSliceField[za0005] = DefinedInt(zb0011)
				}
			}
		case "DefinedStringSliceField":
			var zb0012 uint32
			zb0012, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DefinedStringSliceField) >= int(zb0012) {
				z.DefinedStringSliceField = (z.DefinedStringSliceField)[:zb0012]
			} else {
				z.DefinedStringSliceField = make([]DefinedString, zb0012)
			}
			for za0006 := range z.DefinedStringSliceField {
				{
					var zb0013 string
					zb0013, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
					z.DefinedStringSliceField[za0006] = DefinedString(zb0013)
				}
			}
		case "DefinedFloatSliceField":
			var zb0014 uint32
			zb0014, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DefinedFloatSliceField) >= int(zb0014) {
				z.DefinedFloatSliceField = (z.DefinedFloatSliceField)[:zb0014]
			} else {
				z.DefinedFloatSliceField = make([]DefinedFloat, zb0014)
			}
			for za0007 := range z.DefinedFloatSliceField {
				{
					var zb0015 float64
					zb0015, bts, err = msgp.ReadFloat64Bytes(bts)
					if err != nil {
						return
					}
					z.DefinedFloatSliceField[za0007] = DefinedFloat(zb0015)
				}
			}
		case "DefinedBoolSliceField":
			var zb0016 uint32
			zb0016, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DefinedBoolSliceField) >= int(zb0016) {
				z.DefinedBoolSliceField = (z.DefinedBoolSliceField)[:zb0016]
			} else {
				z.DefinedBoolSliceField = make([]DefinedBool, zb0016)
			}
			for za0008 := range z.DefinedBoolSliceField {
				{
					var zb0017 bool
					zb0017, bts, err = msgp.ReadBoolBytes(bts)
					if err != nil {
						return
					}
					z.DefinedBoolSliceField[za0008] = DefinedBool(zb0017)
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TestType) Msgsize() (s int) {
	s = 3 + 6 + z.Model.Msgsize() + 9 + msgp.IntSize + 12 + msgp.StringPrefixSize + len(z.StringField) + 11 + msgp.Float64Size + 10 + msgp.BoolSize + 14 + msgp.ArrayHeaderSize + (len(z.IntSliceField) * (msgp.IntSize)) + 17 + msgp.ArrayHeaderSize
	for za0002 := range z.StringSliceField {
		s += msgp.StringPrefixSize + len(z.StringSliceField[za0002])
	}
	s += 16 + msgp.ArrayHeaderSize + (len(z.FloatSliceField) * (msgp.Float64Size)) + 15 + msgp.ArrayHeaderSize + (len(z.BoolSliceField) * (msgp.BoolSize)) + 16 + msgp.IntSize + 19 + msgp.StringPrefixSize + len(string(z.DefinedStringField)) + 18 + msgp.Float64Size + 17 + msgp.BoolSize + 21 + msgp.ArrayHeaderSize + (len(z.DefinedIntSliceField) * (msgp.IntSize)) + 24 + msgp.ArrayHeaderSize
	for za0006 := range z.DefinedStringSliceField {
		s += msgp.StringPrefixSize + len(string(z.DefinedStringSliceField[za0006]))
	}
	s += 23 + msgp.ArrayHeaderSize + (len(z.DefinedFloatSliceField) * (msgp.Float64Size)) + 22 + msgp.ArrayHeaderSize + (len(z.DefinedBoolSliceField) * (msgp.BoolSize))
	return
}
