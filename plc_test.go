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
		info, err := client.GetAgBlockInfo(65, 12)
		t.Log(err)
		t.Logf("BlkType: %v", info.BlkType)
		t.Logf("BlkNumber: %v", info.BlkNumber)
		t.Logf("BlkLang: %v", info.BlkLang)
		t.Logf("BlkFlags: %v", info.BlkFlags)
		t.Logf("MC7Size: %v", info.MC7Size)
		t.Logf("LoadSize: %v", info.LoadSize)
		t.Logf("LocalData: %v", info.LocalData)
		t.Logf("SBBLength: %v", info.SBBLength)
		t.Logf("CheckSum: %v", info.CheckSum)
		t.Logf("Version: %v", info.Version)
		t.Logf("CodeDate: %v", info.CodeDate)
		t.Logf("IntfDate: %v", info.IntfDate)
		t.Logf("Author: %v", info.Author)
		t.Logf("Family: %v", info.Family)
		t.Logf("Header: %v", info.Header)
	} else {
		t.Fatal(err)
	}
}
