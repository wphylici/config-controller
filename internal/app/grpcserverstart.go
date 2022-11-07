package app

import (
	"github.com/wphylici/contest-cloud/internal/transport/grpc/server"
	"net"
)

func StartGRPCServer(config *server.Config) error {
	s := server.NewGRPCServer()

	l, err := net.Listen(config.Network, config.BindAddr)
	if err != nil {
		return err
	}

	if err = s.Serve(l); err != nil {
		return err
	}

	return nil
}
