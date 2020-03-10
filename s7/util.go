package s7

import (
	"fmt"
	"sort"

	pb "github.com/thinkontrolsy/goplc/s7/plc_api"
)

type S7AGPointer struct {
	Start    int
	End      int
	DBNumber int
	Buffer   []byte
	Tags     []*pb.Tag
}

func (t S7AGPointer) String() string {
	return fmt.Sprintf("DB: %d | Start: %d | End: %d \n Tags: %v", t.DBNumber, t.Start, t.End, t.Tags)
}

func (t *S7AGPointer) HasBoolTag() bool {
	for _, tag := range t.Tags {
		if tag.Dt == "Bool" {
			return true
		}
	}
	return false
}

func (ag *S7AGPointer) FillBuffer() {
	for _, tag := range ag.Tags {
		if address, err := tag.GetArea(); err == nil {
			buffer := tag.GetValueBuffer()
			for i, b := range buffer {
				ag.Buffer[address.Start-ag.Start+i] = b
			}
		}
	}
}

func (ag *S7AGPointer) ReadBuffer() {
	for _, tag := range ag.Tags {
		if address, err := tag.GetArea(); err == nil {
			s := address.Start - ag.Start
			d := s + address.Amount
			tag.SetValue(ag.Buffer[s:d])
		}
	}
}

func generateAGMap(tags []*pb.Tag) map[string]*S7AGPointer {
	m := make(map[string]*S7AGPointer)
	for _, tag := range tags {
		if address, err := tag.GetArea(); err == nil {
			d, ok := m[address.Area]
			if ok {
				if address.Start < d.Start {
					d.Start = address.Start
				}
				if address.Start+address.Amount > d.End {
					d.End = address.Start + address.Amount
				}
				d.Tags = append(d.Tags, tag)
			} else {
				m[address.Area] = &S7AGPointer{
					Start:    address.Start,
					End:      address.Start + address.Amount,
					DBNumber: address.DBNumber,
					Tags:     []*pb.Tag{tag},
				}
			}
		}

	}
	return m
}

func generateAGGroupMap(tags []*pb.Tag) map[string][]*S7AGPointer {
	groupList := make(map[string][]*pb.Tag)
	for _, tag := range tags {
		if address, err := tag.GetArea(); err == nil {
			groupList[address.Area] = append(groupList[address.Area], tag)
		}
	}
	m := make(map[string][]*S7AGPointer)
	for area, tags := range groupList {
		sort.Slice(tags, func(i, j int) bool {
			add_i, _ := tags[i].GetArea()
			add_j, _ := tags[j].GetArea()
			return add_i.Start < add_j.Start
		})
		var ags []*S7AGPointer
		for _, tag := range tags {
			address, _ := tag.GetArea()
			if len(ags) == 0 {
				ags = append(ags, &S7AGPointer{
					Start:    address.Start,
					End:      address.Start + address.Amount,
					DBNumber: address.DBNumber,
					Tags:     []*pb.Tag{tag},
				})
			} else {
				ag := ags[len(ags)-1]
				if ag.End > address.Start {
					ags = append(ags, &S7AGPointer{
						Start:    address.Start,
						End:      address.Start + address.Amount,
						DBNumber: address.DBNumber,
						Tags:     []*pb.Tag{tag},
					})
				} else {
					if address.Start+address.Amount > ag.End {
						ag.End = address.Start + address.Amount
					}
					ag.Tags = append(ag.Tags, tag)
				}
			}
		}
		m[area] = ags
	}
	return m
}
