package tormenta

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Line) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "ProductCode":
			z.ProductCode, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Qty":
			z.Qty, err = dc.ReadInt()
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
func (z Line) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ProductCode"
	err = en.Append(0x82, 0xab, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.ProductCode)
	if err != nil {
		return
	}
	// write "Qty"
	err = en.Append(0xa3, 0x51, 0x74, 0x79)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Qty)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Line) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ProductCode"
	o = append(o, 0x82, 0xab, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65)
	o = msgp.AppendString(o, z.ProductCode)
	// string "Qty"
	o = append(o, 0xa3, 0x51, 0x74, 0x79)
	o = msgp.AppendInt(o, z.Qty)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Line) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ProductCode":
			z.ProductCode, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Qty":
			z.Qty, bts, err = msgp.ReadIntBytes(bts)
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
func (z Line) Msgsize() (s int) {
	s = 1 + 12 + msgp.StringPrefixSize + len(z.ProductCode) + 4 + msgp.IntSize
	return
}

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
			z.Customer, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "Items":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Items) >= int(zb0002) {
				z.Items = (z.Items)[:zb0002]
			} else {
				z.Items = make([]Line, zb0002)
			}
			for za0001 := range z.Items {
				var zb0003 uint32
				zb0003, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "ProductCode":
						z.Items[za0001].ProductCode, err = dc.ReadString()
						if err != nil {
							return
						}
					case "Qty":
						z.Items[za0001].Qty, err = dc.ReadInt()
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
	// map header, size 3
	// write "Model"
	err = en.Append(0x83, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
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
	err = en.WriteInt(z.Customer)
	if err != nil {
		return
	}
	// write "Items"
	err = en.Append(0xa5, 0x49, 0x74, 0x65, 0x6d, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Items)))
	if err != nil {
		return
	}
	for za0001 := range z.Items {
		// map header, size 2
		// write "ProductCode"
		err = en.Append(0x82, 0xab, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65)
		if err != nil {
			return
		}
		err = en.WriteString(z.Items[za0001].ProductCode)
		if err != nil {
			return
		}
		// write "Qty"
		err = en.Append(0xa3, 0x51, 0x74, 0x79)
		if err != nil {
			return
		}
		err = en.WriteInt(z.Items[za0001].Qty)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Order) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Model"
	o = append(o, 0x83, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
	o, err = z.Model.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "Customer"
	o = append(o, 0xa8, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72)
	o = msgp.AppendInt(o, z.Customer)
	// string "Items"
	o = append(o, 0xa5, 0x49, 0x74, 0x65, 0x6d, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Items)))
	for za0001 := range z.Items {
		// map header, size 2
		// string "ProductCode"
		o = append(o, 0x82, 0xab, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65)
		o = msgp.AppendString(o, z.Items[za0001].ProductCode)
		// string "Qty"
		o = append(o, 0xa3, 0x51, 0x74, 0x79)
		o = msgp.AppendInt(o, z.Items[za0001].Qty)
	}
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
			z.Customer, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "Items":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Items) >= int(zb0002) {
				z.Items = (z.Items)[:zb0002]
			} else {
				z.Items = make([]Line, zb0002)
			}
			for za0001 := range z.Items {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "ProductCode":
						z.Items[za0001].ProductCode, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "Qty":
						z.Items[za0001].Qty, bts, err = msgp.ReadIntBytes(bts)
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
	s = 1 + 6 + z.Model.Msgsize() + 9 + msgp.IntSize + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Items {
		s += 1 + 12 + msgp.StringPrefixSize + len(z.Items[za0001].ProductCode) + 4 + msgp.IntSize
	}
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
		case "Description":
			z.Description, err = dc.ReadString()
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
	// write "Code"
	err = en.Append(0x85, 0xa4, 0x43, 0x6f, 0x64, 0x65)
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
	// write "Description"
	err = en.Append(0xab, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Description)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Product) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "Code"
	o = append(o, 0x85, 0xa4, 0x43, 0x6f, 0x64, 0x65)
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
	// string "Description"
	o = append(o, 0xab, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Description)
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
		case "Description":
			z.Description, bts, err = msgp.ReadStringBytes(bts)
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
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Code) + 5 + msgp.StringPrefixSize + len(z.Name) + 6 + msgp.Float32Size + 14 + msgp.IntSize + 12 + msgp.StringPrefixSize + len(z.Description)
	return
}
