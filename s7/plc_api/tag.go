package plc_api

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	duration "github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	ptypes "github.com/golang/protobuf/ptypes"
	gos7 "github.com/robinson/gos7"
)

const (
	DT_REG  = `^String\[(\d+)\]$`
	ADD_REG = `^(M|I|Q|(?:DB(\d+)))P(\d+)(?:\.([0-7]))?$`
)

var DT = map[string]int{
	"Bool": 1,
	"Byte": 1,
	"Char": 1,

	"Word":  2,
	"DWord": 4,
	"LWord": 8,

	"SInt":  1,
	"USInt": 1,

	"Int":  2,
	"UInt": 2,

	"DInt":  4,
	"UDInt": 4,

	"LInt":  8,
	"ULInt": 8,

	"Real":  4,
	"LReal": 8,

	"DTL": 12,

	"Date":          2,
	"Date_And_Time": 8,
	"LDT":           8,

	"LTime":        8,
	"LTime_Of_Day": 8,

	"S5Time":      2,
	"Time":        4,
	"Time_Of_Day": 4,

	"String": 256,
}

type TagAddress struct {
	Area     string
	DBNumber int
	Start    int
	Bit      uint
	Amount   int
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func decodeBcd(b byte) int {
	return int(((b >> 4) * 10) + (b & 0x0F))
}

func encodeBcd(value int) byte {
	return byte(((value / 10) << 4) | (value % 10))
}

func (tag *Tag) GetLength() int {
	if length, ok := DT[tag.GetDt()]; ok {
		return length
	} else {
		reg, _ := regexp.Compile(DT_REG)
		match := reg.FindStringSubmatch(tag.GetDt())
		if match != nil {
			length, _ = strconv.Atoi(match[1])
			return length + 2
		}
		return 0
	}
}

func (tag *Tag) GetArea() (*TagAddress, error) {
	amount := tag.GetLength()
	if amount == 0 {
		return nil, fmt.Errorf("Datatype illegal")
	}
	reg, _ := regexp.Compile(ADD_REG)
	match := reg.FindStringSubmatch(tag.GetAddress())

	if match == nil {
		return nil, fmt.Errorf("Address illegal")
	} else {
		dbNum, _ := strconv.Atoi(match[2])
		start, _ := strconv.Atoi(match[3])
		bit, _ := strconv.Atoi(match[4])
		return &TagAddress{
			Area:     match[1],
			Amount:   amount,
			DBNumber: dbNum,
			Start:    start,
			Bit:      uint(bit),
		}, nil
	}
}

func (tag *Tag) FillBuffer(b byte) []byte {
	var helper gos7.Helper
	buffer := make([]byte, tag.GetLength())
	switch tag.GetDt() {
	case "Bool":
		{
			address, _ := tag.GetArea()
			v := tag.GetValueBool()
			b_ := helper.SetBoolAt(b, address.Bit, v)
			buffer[0] = b_
		}
	case "Byte", "Word", "DWord", "LWord":
		{
			v := tag.GetValueBytes()
			l := Min(len(v), len(buffer))
			copy(buffer[:l], v[:l])
		}
	case "Char":
		{
			v := tag.GetValueString()
			if len(v) > 0 {
				buffer[0] = byte(v[0])
			}
		}
	case "SInt":
		{
			v := int8(tag.GetValueInteger())
			helper.SetValueAt(buffer, 0, v)
		}
	case "USInt":
		{
			v := uint8(tag.GetValueInteger())
			helper.SetValueAt(buffer, 0, v)
		}
	case "Int":
		{
			v := int16(tag.GetValueInteger())
			helper.SetValueAt(buffer, 0, v)
		}
	case "UInt":
		{
			v := uint16(tag.GetValueInteger())
			helper.SetValueAt(buffer, 0, v)
		}
	case "DInt":
		{
			v := int32(tag.GetValueInteger())
			helper.SetValueAt(buffer, 0, v)
		}
	case "UDInt":
		{
			v := uint32(tag.GetValueInteger())
			helper.SetValueAt(buffer, 0, v)
		}
	case "LInt":
		{
			v := tag.GetValueInteger()
			helper.SetValueAt(buffer, 0, v)
		}
	case "ULInt":
		{
			v := tag.GetValueUinteger()
			helper.SetValueAt(buffer, 0, v)
		}
	case "Real":
		{
			v := float32(tag.GetValueDouble())
			helper.SetValueAt(buffer, 0, v)
		}
	case "LReal":
		{
			v := float64(tag.GetValueDouble())
			helper.SetValueAt(buffer, 0, v)
		}
	case "DTL":
		{
			v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			helper.SetDTLAt(buffer, 0, v)
			// v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			// year := uint16(v.Year())
			// helper.SetValueAt(buffer, 0, year)
			// buffer[2] = byte(v.Month())
			// buffer[3] = byte(v.Day())
			// buffer[4] = byte(v.Weekday())
			// buffer[5] = byte(v.Hour())
			// buffer[6] = byte(v.Minute())
			// buffer[7] = byte(v.Second())
			// nanos := int32(v.Nanosecond())
			// helper.SetValueAt(buffer, 8, nanos)
		}
	case "Date":
		{
			v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			helper.SetDateAt(buffer, 0, v)
			// v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			// initDate := time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
			// hours := v.Sub(initDate).Hours()
			// days := int16(hours / 24)
			// helper.SetValueAt(buffer, 0, days)
		}
	case "Date_And_Time":
		{
			v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			helper.SetDateTimeAt(buffer, 0, v)
			// v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			// year := v.Year()
			// if year >= 2000 {
			// 	year -= 2000
			// } else {
			// 	year -= 1900
			// }
			// buffer[0] = encodeBcd(year)
			// buffer[1] = encodeBcd(int(v.Month()))
			// buffer[2] = encodeBcd(v.Day())
			// buffer[3] = encodeBcd(v.Hour())
			// buffer[4] = encodeBcd(v.Minute())
			// buffer[5] = encodeBcd(v.Second())
			// buffer[6] = encodeBcd(v.Nanosecond() / 1000000 / 10)
			// buffer[7] = (encodeBcd(v.Nanosecond()/1000000%10) << 4) | encodeBcd(int(v.Weekday()))
		}
	case "LDT":
		{
			v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			helper.SetLDTAt(buffer, 0, v)
			// v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			// helper.SetValueAt(buffer, 0, v.UnixNano())
		}
	case "LTime":
		{
			v, _ := ptypes.Duration(tag.GetValueDuration())
			ns := v.Nanoseconds()
			helper.SetValueAt(buffer, 0, ns)
		}
	case "LTime_Of_Day":
		{
			v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			helper.SetLTODAt(buffer, 0, v)
			// t, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			// v := int64((t.Hour()*3600 + t.Minute()*60 + t.Second()) * 1000000000)
			// helper.SetValueAt(buffer, 0, v)
		}
	case "S5Time":
		{
			v, _ := ptypes.Duration(tag.GetValueDuration())
			helper.SetS5TimeAt(buffer, 0, v)
			// ms := v.Milliseconds()
			// switch {
			// case ms < 9990:
			// 	buffer[1] = encodeBcd(int(ms) / 10 % 100)
			// 	buffer[0] = encodeBcd(int(ms)/10/100) &^ 0b11110000
			// case ms > 100 && ms < 99900:
			// 	buffer[1] = encodeBcd(int(ms) / 100 % 100)
			// 	buffer[0] = encodeBcd(int(ms)/100/100)&^0b11100000 | 0b00010000
			// case ms > 1000 && ms < 999000:
			// 	buffer[1] = encodeBcd(int(ms) / 1000 % 100)
			// 	buffer[0] = encodeBcd(int(ms)/1000/100)&^0b11010000 | 0b00100000
			// case ms > 10000 && ms < 9990000:
			// 	buffer[1] = encodeBcd(int(ms) / 10000 % 100)
			// 	buffer[0] = encodeBcd(int(ms)/10000/100)&^0b11000000 | 0b00110000
			// }
		}
	case "Time":
		{
			v, _ := ptypes.Duration(tag.GetValueDuration())
			ms := int32(v.Milliseconds())
			helper.SetValueAt(buffer, 0, ms)
		}
	case "Time_Of_Day":
		{
			v, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			helper.SetTODAt(buffer, 0, v)
			// t, _ := ptypes.Timestamp(tag.GetValueTimestamp())
			// v := int32((t.Hour()*3600 + t.Minute()*60 + t.Second()) * 1000)
			// helper.SetValueAt(buffer, 0, v)
		}
	// String or String[n]
	default:
		{
			v := tag.GetValueString()
			l := Min(len(v), tag.GetLength()-2)
			log.Printf("tag address: %v, l: %v, length: %v", tag.GetAddress(), l, tag.GetLength())
			buffer[0] = byte(tag.GetLength() - 2)
			buffer[1] = byte(l)
			copy(buffer[2:2+l], []byte(v)[:l])
		}
	}
	return buffer
}

func (tag *Tag) SetTagValue(buffer []byte) {
	var helper gos7.Helper
	switch tag.GetDt() {
	case "Bool":
		{
			address, _ := tag.GetArea()
			tag.Value = &Tag_ValueBool{ValueBool: helper.GetBoolAt(buffer[0], address.Bit)}
		}
	case "Byte", "Word", "DWord", "LWord":
		{
			v := make([]byte, tag.GetLength())
			l := Min(len(v), len(buffer))
			copy(v[:l], buffer[:l])
			tag.Value = &Tag_ValueBytes{ValueBytes: v}
		}
	case "Char":
		{
			tag.Value = &Tag_ValueString{ValueString: string(buffer[0])}
		}
	case "SInt":
		{
			var v int8
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: int64(v)}
		}
	case "USInt":
		{
			var v uint8
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: int64(v)}
		}
	case "Int":
		{
			var v int16
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: int64(v)}
		}
	case "UInt":
		{
			var v uint16
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: int64(v)}
		}
	case "DInt":
		{
			var v int32
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: int64(v)}
		}
	case "UDInt":
		{
			var v uint32
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: int64(v)}
		}
	case "LInt":
		{
			var v int64
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueInteger{ValueInteger: v}
		}
	case "ULInt":
		{
			var v uint64
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueUinteger{ValueUinteger: v}
		}
	case "Real":
		{
			var v float32
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueDouble{ValueDouble: float64(v)}
		}
	case "LReal":
		{
			var v float64
			helper.GetValueAt(buffer, 0, &v)
			tag.Value = &Tag_ValueDouble{ValueDouble: v}
		}
	case "DTL":
		{
			v, _ := ptypes.TimestampProto(helper.GetDTLAt(buffer, 0))
			tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
			// var year uint16
			// var nanos int32
			// helper.GetValueAt(buffer, 0, &year)
			// helper.GetValueAt(buffer, 8, &nanos)
			// v, _ := ptypes.TimestampProto(time.Date(int(year), time.Month(int(buffer[2])), int(buffer[3]), int(buffer[5]), int(buffer[6]), int(buffer[7]), int(nanos), time.UTC))
			// tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
		}
	case "Date":
		{
			v, _ := ptypes.TimestampProto(helper.GetDateAt(buffer, 0))
			tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
			// initDate := time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
			// var span int16
			// helper.GetValueAt(buffer, 0, &span)
			// t := initDate.AddDate(0, 0, int(span))
			// v, _ := ptypes.TimestampProto(t)
			// tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
		}
	case "Date_And_Time":
		{
			v, _ := ptypes.TimestampProto(helper.GetDateTimeAt(buffer, 0))
			tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
			// year := decodeBcd(buffer[0])
			// if year >= 90 {
			// 	year += 1900
			// } else {
			// 	year += 2000
			// }
			// month := decodeBcd(buffer[1])
			// day := decodeBcd(buffer[2])
			// hour := decodeBcd(buffer[3])
			// minute := decodeBcd(buffer[4])
			// second := decodeBcd(buffer[5])
			// ms := decodeBcd(buffer[6])*10 + decodeBcd(buffer[7]>>4)
			// v, _ := ptypes.TimestampProto(time.Date(int(year), time.Month(month), day, hour, minute, second, ms*1000000, time.UTC))
			// tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
		}
	case "LDT":
		{
			v, _ := ptypes.TimestampProto(helper.GetLDTAt(buffer, 0))
			tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
			// var nano int64
			// helper.GetValueAt(buffer, 0, &nano)
			// v, _ := ptypes.TimestampProto(time.Date(1970, time.Month(1), 1, 0, 0, 0, int(nano), time.UTC))
			// tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
		}
	case "LTime_Of_Day":
		{
			v, _ := ptypes.TimestampProto(helper.GetLTODAt(buffer, 0))
			tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
			// var nano int64
			// helper.GetValueAt(buffer, 0, &nano)
			// v, _ := ptypes.TimestampProto(time.Date(1970, time.Month(1), 1, 0, 0, 0, int(nano), time.UTC))
			// tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
		}
	case "S5Time":
		{
			// t := decodeBcd(buffer[0]&0b00001111)*100 + decodeBcd(buffer[1])
			// switch buffer[0] & 0b00110000 {
			// case 0b00000000:
			// 	t *= 10
			// case 0b00010000:
			// 	t *= 100
			// case 0b00100000:
			// 	t *= 1000
			// case 0b00110000:
			// 	t *= 10000
			// }
			// d, _ := time.ParseDuration(fmt.Sprintf("%dms", t))
			d := helper.GetS5TimeAt(buffer, 0)
			tag.Value = &Tag_ValueDuration{ValueDuration: ptypes.DurationProto(d)}
		}
	case "Time":
		{
			var ms int32
			helper.GetValueAt(buffer, 0, &ms)
			d, _ := time.ParseDuration(fmt.Sprintf("%dms", ms))
			tag.Value = &Tag_ValueDuration{ValueDuration: ptypes.DurationProto(d)}
		}
	case "LTime":
		{
			var ns int64
			helper.GetValueAt(buffer, 0, &ns)
			d, _ := time.ParseDuration(fmt.Sprintf("%dns", ns))
			tag.Value = &Tag_ValueDuration{ValueDuration: ptypes.DurationProto(d)}
		}
	case "Time_Of_Day":
		{
			v, _ := ptypes.TimestampProto(helper.GetTODAt(buffer, 0))
			tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
			// var ms int32
			// helper.GetValueAt(buffer, 0, &ms)
			// v, _ := ptypes.TimestampProto(time.Date(1970, time.Month(1), 1, 0, 0, 0, int(ms)*1000000, time.UTC))
			// tag.Value = &Tag_ValueTimestamp{ValueTimestamp: v}
		}
	// String or String[n]
	default:
		{
			tag.Value = &Tag_ValueString{ValueString: string(buffer[2:])}
		}
	}
}

func (tag *Tag) GetTagValue() interface{} {
	switch tag.GetDt() {
	case "Bool":
		{
			return tag.GetValueBool()
		}
	case "Byte", "Word", "DWord", "LWord":
		{
			return tag.GetValueBytes()
		}
	case "Char":
		{
			return tag.GetValueString()
		}
	case "SInt", "USInt", "Int", "UInt", "DInt", "UDInt", "LInt", "ULInt":
		{
			return tag.GetValueInteger()
		}
	case "Real", "LReal":
		{
			return tag.GetValueDouble()
		}
	case "DTL", "Date", "Date_And_Time", "LDT", "LTime_Of_Day", "Time_Of_Day":
		{
			return tag.GetValueTimestamp()
		}
	case "LTime", "S5Time", "Time":
		{
			return tag.GetValueDuration()
		}
	// String or String[n]
	default:
		{
			return tag.GetValueString()
		}
	}
}

func (tag *Tag) GetTagValueString() string {
	v := tag.GetTagValue()
	if t, ok := v.(*timestamp.Timestamp); ok {
		tt, _ := ptypes.Timestamp(t)
		return fmt.Sprintf("Address: %s: \t| DT: %s \t| Value: %v\n", tag.GetAddress(), tag.GetDt(), tt)
	}
	if d, ok := v.(*duration.Duration); ok {
		dd, _ := ptypes.Duration(d)
		return fmt.Sprintf("Address: %s: \t| DT: %s \t| Value: %v\n", tag.GetAddress(), tag.GetDt(), dd)
	}
	return fmt.Sprintf("Address: %s: \t| DT: %s \t| Value: %v\n", tag.GetAddress(), tag.GetDt(), v)

}
