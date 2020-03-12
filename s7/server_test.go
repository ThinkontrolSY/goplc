package s7

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/thinkontrolsy/goplc/s7/plc_api"
)

func TestGetCpuInfo(t *testing.T) {
	t.Log("Test CPU Info")
	server := PlcServer{}

	plc := &pb.Plc{
		Host: "10.0.0.230",
		Rack: 0,
		Slot: 1,
	}

	info, err := server.GetCpuInfo(context.Background(), plc)

	t.Log(info)
	t.Log(err)
}

func TestWriteTags(t *testing.T) {
	t.Log("Test Write Tags")
	server := PlcServer{}

	plc := &pb.Plc{
		Host: "10.0.0.230",
		Rack: 0,
		Slot: 1,
	}

	now := time.Now()
	du, _ := time.ParseDuration("1h50ms")
	ts, _ := ptypes.TimestampProto(now)
	dup := ptypes.DurationProto(du)

	t.Log("//////////////////////////////////////")
	t.Log(now)
	t.Log(du)

	tags := []*pb.Tag{
		&pb.Tag{
			Address: "DB2P2",
			Dt:      "Int",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 2345},
		},
		&pb.Tag{
			Address: "DB2P0",
			Dt:      "Int",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 2345},
		},
		&pb.Tag{
			Address: "DB2P4",
			Dt:      "Real",
			Value:   &pb.Tag_ValueDouble{ValueDouble: 234.123},
		},
		&pb.Tag{
			Address: "DB2P9.6",
			Dt:      "Bool",
			Value:   &pb.Tag_ValueBool{ValueBool: true},
		},
		&pb.Tag{
			Address: "DB2P9.2",
			Dt:      "Bool",
			Value:   &pb.Tag_ValueBool{ValueBool: true},
		},
		&pb.Tag{
			Address: "DB2P9.1",
			Dt:      "Bool",
			Value:   &pb.Tag_ValueBool{ValueBool: false},
		},
		&pb.Tag{
			Address: "DB2P9.4",
			Dt:      "Bool",
			Value:   &pb.Tag_ValueBool{ValueBool: false},
		},
		&pb.Tag{
			Address: "DB2P10",
			Dt:      "DInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 234123},
		},
		&pb.Tag{
			Address: "DB2P18",
			Dt:      "DTL",
			Value:   &pb.Tag_ValueTimestamp{ValueTimestamp: ts},
		},
		&pb.Tag{
			Address: "DB2P30",
			Dt:      "Date",
			Value:   &pb.Tag_ValueTimestamp{ValueTimestamp: ts},
		},
		&pb.Tag{
			Address: "DB2P32",
			Dt:      "Date_And_Time",
			Value:   &pb.Tag_ValueTimestamp{ValueTimestamp: ts},
		},
		&pb.Tag{
			Address: "DB2P40",
			Dt:      "LDT",
			Value:   &pb.Tag_ValueTimestamp{ValueTimestamp: ts},
		},
		&pb.Tag{
			Address: "DB2P48",
			Dt:      "LInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: -234123},
		},
		&pb.Tag{
			Address: "DB2P56",
			Dt:      "LReal",
			Value:   &pb.Tag_ValueDouble{ValueDouble: 234.123},
		},
		&pb.Tag{
			Address: "DB2P64",
			Dt:      "LTime",
			Value:   &pb.Tag_ValueDuration{ValueDuration: dup},
		},
		&pb.Tag{
			Address: "DB2P72",
			Dt:      "LTime_Of_Day",
			Value:   &pb.Tag_ValueTimestamp{ValueTimestamp: ts},
		},
		&pb.Tag{
			Address: "DB2P80",
			Dt:      "S5Time",
			Value:   &pb.Tag_ValueDuration{ValueDuration: dup},
		},
		&pb.Tag{
			Address: "DB2P82",
			Dt:      "String[12]",
			Value:   &pb.Tag_ValueString{ValueString: "abcdefghijklmnopqrstuvwxyz"},
		},
		&pb.Tag{
			Address: "DB2P96",
			Dt:      "String",
			Value:   &pb.Tag_ValueString{ValueString: "abcdefghijklmnopqrstuvwxyz"},
		},
		&pb.Tag{
			Address: "DB2P352",
			Dt:      "Time",
			Value:   &pb.Tag_ValueDuration{ValueDuration: dup},
		},
		&pb.Tag{
			Address: "DB2P356",
			Dt:      "Time_Of_Day",
			Value:   &pb.Tag_ValueTimestamp{ValueTimestamp: ts},
		},
		&pb.Tag{
			Address: "DB2P360",
			Dt:      "Byte",
			Value:   &pb.Tag_ValueBytes{ValueBytes: []byte{2, 1, 4, 5, 2, 12, 4, 1, 1, 3, 4, 9, 21, 4, 23, 1, 7}},
		},
		&pb.Tag{
			Address: "DB2P361",
			Dt:      "Char",
			Value:   &pb.Tag_ValueString{ValueString: "vvda"},
		},
		&pb.Tag{
			Address: "DB2P362",
			Dt:      "DWord",
			Value:   &pb.Tag_ValueBytes{ValueBytes: []byte{2, 1, 4, 5, 2, 12, 4, 1, 1, 3, 4, 9, 21, 4, 23, 1, 7}},
		},
		&pb.Tag{
			Address: "DB2P366",
			Dt:      "LWord",
			Value:   &pb.Tag_ValueBytes{ValueBytes: []byte{2, 1, 4, 5, 2, 12, 4, 1, 1, 3, 4, 9, 21, 4, 23, 1, 7}},
		},
		&pb.Tag{
			Address: "DB2P374",
			Dt:      "SInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: -2},
		},
		&pb.Tag{
			Address: "DB2P376",
			Dt:      "UDInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 234123},
		},
		&pb.Tag{
			Address: "DB2P380",
			Dt:      "UInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 234123},
		},
		&pb.Tag{
			Address: "DB2P382",
			Dt:      "ULInt",
			Value:   &pb.Tag_ValueUinteger{ValueUinteger: 234123},
		},
		&pb.Tag{
			Address: "DB2P390",
			Dt:      "USInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 4},
		},
		&pb.Tag{
			Address: "DB2P391",
			Dt:      "USInt",
			Value:   &pb.Tag_ValueInteger{ValueInteger: 234123},
		},
	}

	_, err := server.WriteTags(context.Background(), &pb.RWReq{Plc: plc, Tags: tags})
	t.Logf("err: %v", err)
}

func TestReadTags(t *testing.T) {
	t.Log("Test Read Tags")
	server := PlcServer{}

	plc := &pb.Plc{
		Host: "10.0.0.230",
		Rack: 0,
		Slot: 1,
	}

	tags := []*pb.Tag{
		&pb.Tag{
			Address: "DB2P2",
			Dt:      "Int",
		},
		&pb.Tag{
			Address: "DB2P0",
			Dt:      "Int",
		},
		&pb.Tag{
			Address: "DB2P4",
			Dt:      "Real",
		},
		&pb.Tag{
			Address: "DB2P9.6",
			Dt:      "Bool",
		},
		&pb.Tag{
			Address: "DB2P9.2",
			Dt:      "Bool",
		},
		&pb.Tag{
			Address: "DB2P9.1",
			Dt:      "Bool",
		},
		&pb.Tag{
			Address: "DB2P9.4",
			Dt:      "Bool",
		},
		&pb.Tag{
			Address: "DB2P10",
			Dt:      "DInt",
		},
		&pb.Tag{
			Address: "DB2P18",
			Dt:      "DTL",
		},
		&pb.Tag{
			Address: "DB2P30",
			Dt:      "Date",
		},
		&pb.Tag{
			Address: "DB2P32",
			Dt:      "Date_And_Time",
		},
		&pb.Tag{
			Address: "DB2P40",
			Dt:      "LDT",
		},
		&pb.Tag{
			Address: "DB2P48",
			Dt:      "LInt",
		},
		&pb.Tag{
			Address: "DB2P56",
			Dt:      "LReal",
		},
		&pb.Tag{
			Address: "DB2P64",
			Dt:      "LTime",
		},
		&pb.Tag{
			Address: "DB2P72",
			Dt:      "LTime_Of_Day",
		},
		&pb.Tag{
			Address: "DB2P80",
			Dt:      "S5Time",
		},
		&pb.Tag{
			Address: "DB2P82",
			Dt:      "String[12]",
		},
		&pb.Tag{
			Address: "DB2P96",
			Dt:      "String",
		},
		&pb.Tag{
			Address: "DB2P352",
			Dt:      "Time",
		},
		&pb.Tag{
			Address: "DB2P356",
			Dt:      "Time_Of_Day",
		},
		&pb.Tag{
			Address: "DB2P360",
			Dt:      "Byte",
		},
		&pb.Tag{
			Address: "DB2P361",
			Dt:      "Char",
		},
		&pb.Tag{
			Address: "DB2P362",
			Dt:      "DWord",
		},
		&pb.Tag{
			Address: "DB2P366",
			Dt:      "LWord",
		},
		&pb.Tag{
			Address: "DB2P374",
			Dt:      "SInt",
		},
		&pb.Tag{
			Address: "DB2P376",
			Dt:      "UDInt",
		},
		&pb.Tag{
			Address: "DB2P380",
			Dt:      "UInt",
		},
		&pb.Tag{
			Address: "DB2P382",
			Dt:      "ULInt",
		},
		&pb.Tag{
			Address: "DB2P390",
			Dt:      "USInt",
		},
		&pb.Tag{
			Address: "DB2P391",
			Dt:      "USInt",
		},
	}

	r, err := server.ReadTags(context.Background(), &pb.RWReq{Plc: plc, Tags: tags})
	t.Logf("err: %v", err)

	for _, tag := range r.GetTags() {
		t.Log(tag.GetTagValueString())
	}
}
