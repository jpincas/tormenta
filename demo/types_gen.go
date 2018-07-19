package demo

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
func (z *NoModel) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "SomeData":
			z.SomeData, err = dc.ReadString()
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
func (z NoModel) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "SomeData"
	err = en.Append(0x81, 0xa8, 0x53, 0x6f, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = en.WriteString(z.SomeData)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z NoModel) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "SomeData"
	o = append(o, 0x81, 0xa8, 0x53, 0x6f, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendString(o, z.SomeData)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NoModel) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "SomeData":
			z.SomeData, bts, err = msgp.ReadStringBytes(bts)
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
func (z NoModel) Msgsize() (s int) {
	s = 1 + 9 + msgp.StringPrefixSize + len(z.SomeData)
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
		case "HasShipped":
			z.HasShipped, err = dc.ReadBool()
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
	// map header, size 6
	// write "Model"
	err = en.Append(0x86, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
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
	// write "HasShipped"
	err = en.Append(0xaa, 0x48, 0x61, 0x73, 0x53, 0x68, 0x69, 0x70, 0x70, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.HasShipped)
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
	// map header, size 6
	// string "Model"
	o = append(o, 0x86, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
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
	// string "HasShipped"
	o = append(o, 0xaa, 0x48, 0x61, 0x73, 0x53, 0x68, 0x69, 0x70, 0x70, 0x65, 0x64)
	o = msgp.AppendBool(o, z.HasShipped)
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
		case "HasShipped":
			z.HasShipped, bts, err = msgp.ReadBoolBytes(bts)
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
	s = 1 + 6 + z.Model.Msgsize() + 9 + msgp.StringPrefixSize + len(z.Customer) + 11 + msgp.IntSize + 12 + msgp.Float64Size + 11 + msgp.BoolSize + 6 + msgp.ArrayHeaderSize
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
			z.Price, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "StartingStock":
			z.StartingStock, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "Tags":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Tags) >= int(zb0002) {
				z.Tags = (z.Tags)[:zb0002]
			} else {
				z.Tags = make([]string, zb0002)
			}
			for za0001 := range z.Tags {
				z.Tags[za0001], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "Departments":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Departments) >= int(zb0003) {
				z.Departments = (z.Departments)[:zb0003]
			} else {
				z.Departments = make([]int, zb0003)
			}
			for za0002 := range z.Departments {
				z.Departments[za0002], err = dc.ReadInt()
				if err != nil {
					return
				}
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
	// map header, size 8
	// write "Model"
	err = en.Append(0x88, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
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
	err = en.WriteFloat64(z.Price)
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
	// write "Tags"
	err = en.Append(0xa4, 0x54, 0x61, 0x67, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Tags)))
	if err != nil {
		return
	}
	for za0001 := range z.Tags {
		err = en.WriteString(z.Tags[za0001])
		if err != nil {
			return
		}
	}
	// write "Departments"
	err = en.Append(0xab, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Departments)))
	if err != nil {
		return
	}
	for za0002 := range z.Departments {
		err = en.WriteInt(z.Departments[za0002])
		if err != nil {
			return
		}
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
	// map header, size 8
	// string "Model"
	o = append(o, 0x88, 0xa5, 0x4d, 0x6f, 0x64, 0x65, 0x6c)
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
	o = msgp.AppendFloat64(o, z.Price)
	// string "StartingStock"
	o = append(o, 0xad, 0x53, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x6f, 0x63, 0x6b)
	o = msgp.AppendInt(o, z.StartingStock)
	// string "Tags"
	o = append(o, 0xa4, 0x54, 0x61, 0x67, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Tags)))
	for za0001 := range z.Tags {
		o = msgp.AppendString(o, z.Tags[za0001])
	}
	// string "Departments"
	o = append(o, 0xab, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Departments)))
	for za0002 := range z.Departments {
		o = msgp.AppendInt(o, z.Departments[za0002])
	}
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
			z.Price, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "StartingStock":
			z.StartingStock, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "Tags":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Tags) >= int(zb0002) {
				z.Tags = (z.Tags)[:zb0002]
			} else {
				z.Tags = make([]string, zb0002)
			}
			for za0001 := range z.Tags {
				z.Tags[za0001], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "Departments":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Departments) >= int(zb0003) {
				z.Departments = (z.Departments)[:zb0003]
			} else {
				z.Departments = make([]int, zb0003)
			}
			for za0002 := range z.Departments {
				z.Departments[za0002], bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					return
				}
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
	s = 1 + 6 + z.Model.Msgsize() + 5 + msgp.StringPrefixSize + len(z.Code) + 5 + msgp.StringPrefixSize + len(z.Name) + 6 + msgp.Float64Size + 14 + msgp.IntSize + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Tags {
		s += msgp.StringPrefixSize + len(z.Tags[za0001])
	}
	s += 12 + msgp.ArrayHeaderSize + (len(z.Departments) * (msgp.IntSize)) + 12 + msgp.StringPrefixSize + len(z.Description)
	return
}
