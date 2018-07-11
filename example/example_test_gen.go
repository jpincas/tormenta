package example

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Order) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Customer":
			z.Customer, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Department":
			z.Department, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "ShippingFee":
			z.ShippingFee, err = dc.ReadFloat64()
			if err != nil {
				return
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
func (z *Order) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Model"
	err = en.Append(0x84, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	if err != nil {
		return
	}
	err = z.Model.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Customer"
	err = en.Append(0xa8, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Customer)
	if err != nil {
		return
	}
	// write "Department"
	err = en.Append(0xaa, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Department)
	if err != nil {
		return
	}
	// write "ShippingFee"
	err = en.Append(0xab, 0x53, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x46, 0x65, 0x65)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.ShippingFee)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Order) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Model"
	o = append(o, 0x84, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	o, err = z.Model.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Customer"
	o = append(o, 0xa8, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72)
	o = msgp.AppendString(o, z.Customer)
	// string "Department"
	o = append(o, 0xaa, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74)
	o = msgp.AppendInt(o, z.Department)
	// string "ShippingFee"
	o = append(o, 0xab, 0x53, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x46, 0x65, 0x65)
	o = msgp.AppendFloat64(o, z.ShippingFee)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Order) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Customer":
			z.Customer, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Department":
			z.Department, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "ShippingFee":
			z.ShippingFee, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
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
func (z *Order) Msgsize() (s int) {
	s = 1 + 6 + z.Model.Msgsize() + 9 + msgp.StringPrefixSize + len(z.Customer) + 11 + msgp.IntSize + 12 + msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Product) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "Code":
			z.Code, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Price":
			z.Price, err = dc.ReadFloat32()
			if err != nil {
				return
			}
		case "StartingStock":
			z.StartingStock, err = dc.ReadInt()
			if err != nil {
				return
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
func (z *Product) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "Model"
	err = en.Append(0x85, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	if err != nil {
		return
	}
	err = z.Model.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "Code"
	err = en.Append(0xa4, 0x43, 0x6f, 0x64, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Code)
	if err != nil {
		return
	}
	// write "Name"
	err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "Price"
	err = en.Append(0xa5, 0x50, 0x72, 0x69, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteFloat32(z.Price)
	if err != nil {
		return
	}
	// write "StartingStock"
	err = en.Append(0xad, 0x53, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	err = en.WriteInt(z.StartingStock)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Product) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Model"
	o = append(o, 0x85, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	o, err = z.Model.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Code"
	o = append(o, 0xa4, 0x43, 0x6f, 0x64, 0x65)
	o = msgp.AppendString(o, z.Code)
	// string "Name"
	o = append(o, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "Price"
	o = append(o, 0xa5, 0x50, 0x72, 0x69, 0x63, 0x65)
	o = msgp.AppendFloat32(o, z.Price)
	// string "StartingStock"
	o = append(o, 0xad, 0x53, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x6f, 0x63, 0x6b)
	o = msgp.AppendInt(o, z.StartingStock)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Product) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Code":
			z.Code, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Price":
			z.Price, bts, err = msgp.ReadFloat32Bytes(bts)
			if err != nil {
				return
			}
		case "StartingStock":
			z.StartingStock, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
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
func (z *Product) Msgsize() (s int) {
	s = 1 + 6 + z.Model.Msgsize() + 5 + msgp.StringPrefixSize + len(z.Code) + 5 + msgp.StringPrefixSize + len(z.Name) + 6 + msgp.Float32Size + 14 + msgp.IntSize
	return
}
