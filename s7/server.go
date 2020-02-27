package s7

import (
	"context"
	"time"

	gos7 "github.com/robinson/gos7"
)

type PlcServer struct {
	UnimplementedPlcRWServer
}

func (s *PlcServer) GetCpuInfo(ctx context.Context, req *Plc) (*S7CpuInfo, error) {
	handler := gos7.NewTCPClientHandler(req.GetHost(), int(req.GetRack()), int(req.GetSlot()))
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 5 * time.Second
	defer handler.Close()
	if err := handler.Connect(); err == nil {
		client := gos7.NewClient(handler)
		info, err := client.GetCPUInfo()
		if err == nil {
			return &S7CpuInfo{
				ModuleTypeName: info.ModuleTypeName,
				SerialNumber:   info.SerialNumber,
				AsName:         info.ASName,
				Copyright:      info.Copyright,
				ModuleName:     info.ModuleName,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
func (s *PlcServer) ReadTags(ctx context.Context, req *RWReq) (*RWResult, error) {
	handler := gos7.NewTCPClientHandler(req.GetPlc().GetHost(), int(req.GetPlc().GetRack()), int(req.GetPlc().GetSlot()))
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 5 * time.Second
	defer handler.Close()
	if err := handler.Connect(); err == nil {
		tags := req.GetTags()
		if items, err := tagsConvert(tags); err == nil {
			client := gos7.NewClient(handler)
			for area, ag := range generateAGMap(items) {
				size := ag.End - ag.Start
				ag.Buffer = make([]byte, size)
				var err error
				switch area {
				case "M":
					err = client.AGReadMB(ag.Start, size, ag.Buffer)
				case "I":
					err = client.AGReadAB(ag.Start, size, ag.Buffer)
				case "Q":
					err = client.AGReadEB(ag.Start, size, ag.Buffer)
				default:
					err = client.AGReadDB(ag.DBNumber, ag.Start, size, ag.Buffer)
				}
				if err != nil {
					return nil, err
				}
				ag.ReadBuffer()
			}
			return &RWResult{Tags: tags}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
func (s *PlcServer) WriteTags(ctx context.Context, req *RWReq) (*RWResult, error) {
	handler := gos7.NewTCPClientHandler(req.GetPlc().GetHost(), int(req.GetPlc().GetRack()), int(req.GetPlc().GetSlot()))
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 5 * time.Second
	defer handler.Close()
	if err := handler.Connect(); err == nil {
		tags := req.GetTags()
		if items, err := tagsConvert(tags); err == nil {
			client := gos7.NewClient(handler)
			for area, ags := range generateAGGroupMap(items) {
				for _, ag := range ags {
					size := ag.End - ag.Start
					ag.Buffer = make([]byte, size)
					switch area {
					case "M":
						{
							if ag.HasBoolTag() {
								if err := client.AGReadMB(ag.Start, size, ag.Buffer); err != nil {
									return nil, err
								}
							}
							ag.FillBuffer()
							if err := client.AGWriteMB(ag.Start, size, ag.Buffer); err != nil {
								return nil, err
							}
						}
					case "I":
						{
							if ag.HasBoolTag() {
								if err := client.AGReadAB(ag.Start, size, ag.Buffer); err != nil {
									return nil, err
								}
							}
							ag.FillBuffer()
							if err := client.AGWriteAB(ag.Start, size, ag.Buffer); err != nil {
								return nil, err
							}
						}
					case "Q":
						{
							if ag.HasBoolTag() {
								if err := client.AGReadEB(ag.Start, size, ag.Buffer); err != nil {
									return nil, err
								}
							}
							ag.FillBuffer()
							if err := client.AGWriteEB(ag.Start, size, ag.Buffer); err != nil {
								return nil, err
							}
						}
					default:
						{
							if ag.HasBoolTag() {
								if err := client.AGReadDB(ag.DBNumber, ag.Start, size, ag.Buffer); err != nil {
									return nil, err
								}
							}
							ag.FillBuffer()
							if err := client.AGWriteDB(ag.DBNumber, ag.Start, size, ag.Buffer); err != nil {
								return nil, err
							}
						}
					}
				}
			}
			return &RWResult{Tags: tags}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
