package s7

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	gos7 "github.com/robinson/gos7"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	AddressRegStr string = `^(M|I|Q|(?:DB(\d+)))(W|D|X)(\d+)(?:\.([0-7]))?$`
)

type S7TagAddress struct {
	Area     string
	DBNumber int
	Start    int
	Bit      uint
	Amount   int
	RawTag   *Tag
}

func (t S7TagAddress) String() string {
	return fmt.Sprintf("Area: %s | DB: %d | Start: %d | Bit: %d | Amount: %d", t.Area, t.DBNumber, t.Start, t.Bit, t.Amount)
}

type S7AGPointer struct {
	Start    int
	End      int
	DBNumber int
	Buffer   []byte
	Tags     []*S7TagAddress
}

func (t S7AGPointer) String() string {
	return fmt.Sprintf("DB: %d | Start: %d | End: %d \n Tags: %v", t.DBNumber, t.Start, t.End, t.Tags)
}

func (t *S7AGPointer) HasBoolTag() bool {
	for _, tag := range t.Tags {
		if tag.RawTag.Dt == DataType_BOOL {
			return true
		}
	}
	return false
}

func (ag *S7AGPointer) FillBuffer() {
	var helper gos7.Helper
	for _, item := range ag.Tags {
		tag := item.RawTag
		switch tag.Dt {
		case DataType_BOOL:
			{
				v := tag.GetValueBool()
				b := helper.SetBoolAt(ag.Buffer[item.Start-ag.Start], item.Bit, v)
				ag.Buffer[item.Start-ag.Start] = b
			}
		case DataType_DINT:
			{
				v := tag.GetValueInt()
				helper.SetValueAt(ag.Buffer, item.Start-ag.Start, v)
			}
		case DataType_INT:
			{
				v := int16(tag.GetValueInt())
				helper.SetValueAt(ag.Buffer, item.Start-ag.Start, v)
			}
		case DataType_REAL:
			{
				v := tag.GetValueFloat()
				helper.SetValueAt(ag.Buffer, item.Start-ag.Start, v)
			}
		}
	}
}

func (ag *S7AGPointer) ReadBuffer() {
	var helper gos7.Helper
	for _, item := range ag.Tags {
		tag := item.RawTag
		switch tag.Dt {
		case DataType_BOOL:
			{
				tag.Value = &Tag_ValueBool{ValueBool: helper.GetBoolAt(ag.Buffer[item.Start-ag.Start], item.Bit)}
			}
		case DataType_DINT:
			{
				var r int32
				helper.GetValueAt(ag.Buffer, item.Start-ag.Start, &r)
				tag.Value = &Tag_ValueInt{ValueInt: r}
			}
		case DataType_INT:
			{
				var r int16
				helper.GetValueAt(ag.Buffer, item.Start-ag.Start, &r)
				tag.Value = &Tag_ValueInt{ValueInt: int32(r)}
			}
		case DataType_REAL:
			{
				var r float32
				helper.GetValueAt(ag.Buffer, item.Start-ag.Start, &r)
				tag.Value = &Tag_ValueFloat{ValueFloat: r}
			}
		}
	}
}

func tagsConvert(tags []*Tag) ([]*S7TagAddress, error) {
	addReg, _ := regexp.Compile(AddressRegStr)
	var items []*S7TagAddress
	for _, tag := range tags {
		match := addReg.FindStringSubmatch(tag.GetAddress())
		if match == nil {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Address format error: %v", tag.GetAddress()))
		}
		var start, amount int
		bit, _ := strconv.Atoi(match[5])
		dbNum, _ := strconv.Atoi(match[2])
		switch match[3] {
		case "W":
			amount = 2
		case "D":
			amount = 4
		default:
			amount = 1
		}
		start, _ = strconv.Atoi(match[4])
		items = append(items, &S7TagAddress{
			Area:     match[1],
			DBNumber: dbNum,
			Start:    start,
			Amount:   amount,
			Bit:      uint(bit),
			RawTag:   tag,
		})
	}
	return items, nil
}

func generateAGMap(tags []*S7TagAddress) map[string]*S7AGPointer {
	m := make(map[string]*S7AGPointer)
	for _, tag := range tags {
		d, ok := m[tag.Area]
		if ok {
			if tag.Start < d.Start {
				d.Start = tag.Start
			}
			if tag.Start+tag.Amount > d.End {
				d.End = tag.Start + tag.Amount
			}
			d.Tags = append(d.Tags, tag)
		} else {
			m[tag.Area] = &S7AGPointer{
				Start:    tag.Start,
				End:      tag.Start + tag.Amount,
				DBNumber: tag.DBNumber,
				Tags:     []*S7TagAddress{tag},
			}
		}
	}
	return m
}

func generateAGGroupMap(tags []*S7TagAddress) map[string][]*S7AGPointer {
	groupList := make(map[string][]*S7TagAddress)
	for _, tag := range tags {
		groupList[tag.Area] = append(groupList[tag.Area], tag)
	}
	m := make(map[string][]*S7AGPointer)
	for area, tags := range groupList {
		sort.Slice(tags, func(i, j int) bool { return tags[i].Start < tags[j].Start })
		var ags []*S7AGPointer
		for _, tag := range tags {
			if len(ags) == 0 {
				ags = append(ags, &S7AGPointer{
					Start:    tag.Start,
					End:      tag.Start + tag.Amount,
					DBNumber: tag.DBNumber,
					Tags:     []*S7TagAddress{tag},
				})
			} else {
				ag := ags[len(ags)-1]
				if ag.End > tag.Start {
					ags = append(ags, &S7AGPointer{
						Start:    tag.Start,
						End:      tag.Start + tag.Amount,
						DBNumber: tag.DBNumber,
						Tags:     []*S7TagAddress{tag},
					})
				} else {
					if tag.Start+tag.Amount > ag.End {
						ag.End = tag.Start + tag.Amount
					}
					ag.Tags = append(ag.Tags, tag)
				}
			}
		}
		m[area] = ags
	}
	return m
}
