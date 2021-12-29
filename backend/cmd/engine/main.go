package main

import (
	"log"
	"net"
	"nft-engine/internal/engine"
	"nft-engine/pkg/proto"
	"time"

	"google.golang.org/grpc"
)

// Retrieve user password and load keystore

func ConfigPath() string {
	return ""
}

func main() {
	log.Println("Starting engine...")
	// engine.LoadConfig("TEST_ID")
	lis, err := net.Listen("tcp", "localhost:17529")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*1024), grpc.MaxSendMsgSize(1024*1024*1024))

	proto.RegisterEngineHandlerServer(grpcServer, EngineHandler{})
	proto.RegisterVaultHandlerServer(grpcServer, VaultHandler{})

	go func() {
		time.Sleep(time.Second)
		engine.Notify(&proto.Notification{
			Message: "Ready",
			Level:   proto.StatusLevel_Log,
		})
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
