package main

import (
	"testing"
	"time"

	// "github.com/golang/protobuf/ptypes"

	gos7 "github.com/robinson/gos7"
)

func TestGetCpuInfo(t *testing.T) {
	handler := gos7.NewTCPClientHandler("10.0.0.230", 0, 1)
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 5 * time.Second
	defer handler.Close()
	if err := handler.Connect(); err == nil {
		client := gos7.NewClient(handler)
		info, err := client.GetCPUInfo()
		t.Log(err)
		t.Log(info)
	} else {
		t.Fatal(err)
	}
}

func TestGetBlockInfo(t *testing.T) {
	handler := gos7.NewTCPClientHandler("10.0.0.230", 0, 1)
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 5 * time.Second
	defer handler.Close()
	if err := handler.Connect(); err == nil {
		client := gos7.NewClient(handler)
		info, err := client.GetAgBlockInfo(65, 6)
		t.Logf("%+v", info)
		t.Log(err)

		list, err := client.PGListBlocks()
		t.Logf("%+v", list)
		t.Log(err)
	} else {
		t.Fatal(err)
	}
}
