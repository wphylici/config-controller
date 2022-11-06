package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/wphylici/contest-cloud/internal/app"
	"github.com/wphylici/contest-cloud/internal/database"
	"github.com/wphylici/contest-cloud/internal/transport/grpc/server"
	"log"
	"os"
)

func main() {
	var gRPCConfigPath string
	var postgreSQLConfigPath string

	const databaseURLEnv = "DATABASE_URL"

	flag.StringVar(&gRPCConfigPath, "grpc_conf", "configs/grpc_server_config.toml", "path to gRPC server config file")
	flag.StringVar(&postgreSQLConfigPath, "postgresql_conf", "configs/postgresql_config.toml", "path to PostgreSQL config file")
	flag.Parse()

	configPostgreSQL := database.NewConfig()
	if app.IsRunningInDockerContainer() {
		configPostgreSQL.DatabaseURL = os.Getenv(databaseURLEnv)
	} else {
		if _, err := toml.DecodeFile(postgreSQLConfigPath, &configPostgreSQL); err != nil {
			log.Fatal(err)
		}
	}

	if err := app.StartPostgreSQL(configPostgreSQL); err != nil {
		log.Fatal(err)
	}

	configGRPCServer := server.NewConfig()
	if _, err := toml.DecodeFile(gRPCConfigPath, &configGRPCServer); err != nil {
		log.Fatal(err)
	}
	if err := app.StartGRPCServer(configGRPCServer); err != nil {
		log.Fatal(err)
	}
}
