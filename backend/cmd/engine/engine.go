package main

import (
	"context"
	"nft-engine/internal/engine"
	"nft-engine/pkg/proto"
	"time"

	"github.com/denisbrodbeck/machineid"
)

const VERSION = "0.0.0"

type EngineHandler struct {
	proto.UnimplementedEngineHandlerServer
}

func (EngineHandler) Init(context.Context, *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func (EngineHandler) Notify(e *proto.Empty, stream proto.EngineHandler_NotifyServer) error {
	for notification := range engine.StatusChan {
		stream.SendMsg(notification)
	}

	return nil
}

func (EngineHandler) Listen(e *proto.Empty, stream proto.EngineHandler_ListenServer) error {
	for range time.Tick(time.Millisecond * 500) {
		hwid, err := machineid.ProtectedID("pawsbackend")
		if err != nil {
			continue
		}

		stream.SendMsg(&proto.EngineStatus{
			Connected: true,
			Version:   VERSION,
			Hwid:      hwid,
		})
	}

	return nil
}
