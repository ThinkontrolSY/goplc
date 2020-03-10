package plc_api

import (
	"fmt"
	"regexp"
	"strconv"
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

func (tag *Tag) GetLength() int {
	if length, ok := DT[tag.GetDt()]; ok {
		return length
	} else {
		reg, _ := regexp.Compile(DT_REG)
		match := reg.FindStringSubmatch(tag.GetDt())
		if match != nil {
			length, _ = strconv.Atoi(match[1])
			return length
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

func (tag *Tag) GetValueBuffer() []byte {
	buffer := make([]byte, tag.GetLength())
	return buffer
}

func (tag *Tag) SetValue(buffer []byte) {

}
