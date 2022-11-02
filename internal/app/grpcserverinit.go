package app

import (
	"github.com/wphylici/contest-cloud/internal/transport/grpc/server"
	"log"
	"net"
)

func StartGRPCServer(config *server.Config) {
	s := server.NewGRPCServer()

	l, err := net.Listen(config.Network, config.BindAddr)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
