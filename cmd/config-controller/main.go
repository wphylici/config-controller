package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/wphylici/contest-cloud/internal/app"
	"github.com/wphylici/contest-cloud/internal/transport/grpc/server"
	"log"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "conf", "configs/grpc_server_config.toml", "path to gRPC server config file")
	flag.Parse()

	configGRPCServer := server.NewConfig()
	_, err := toml.DecodeFile(configPath, configGRPCServer)
	if err != nil {
		log.Fatal(err)
	}

	app.StartGRPCServer(configGRPCServer)
}
