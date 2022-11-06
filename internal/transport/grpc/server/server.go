package server

import (
	"context"
	"encoding/json"
	"github.com/wphylici/contest-cloud/internal/database"
	"github.com/wphylici/contest-cloud/internal/models"
	"github.com/wphylici/contest-cloud/internal/transport/grpc/pb"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	pb.UnimplementedConfigControllerServer
}

func NewGRPCServer() *grpc.Server {
	srv := gRPCServer{}

	s := grpc.NewServer()
	pb.RegisterConfigControllerServer(s, &srv)
	return s
}

func (s *gRPCServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {

	serviceConfig := &models.ServiceConfig{}
	err := json.Unmarshal([]byte(req.ConfData), &serviceConfig)
	if err != nil {
		return nil, err
	}

	screp := database.Psql.ServiceConfig()
	serviceConfig, err = screp.Create(serviceConfig)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Resp: "Success"}, nil
}

func (s *gRPCServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {

	screp := database.Psql.ServiceConfig()
	serviceConfig, err := screp.Read(&models.ServiceConfig{
		Service: req.ServiceName,
		Version: req.Version,
	})
	if err != nil {
		return nil, err
	}

	configData, err := json.Marshal(serviceConfig.Data)
	if err != nil {
		return nil, err
	}

	return &pb.ReadResponse{Resp: "Success", ConfData: string(configData)}, nil
}

func (s *gRPCServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	serviceConfig := &models.ServiceConfig{}
	err := json.Unmarshal([]byte(req.ConfData), &serviceConfig)
	if err != nil {
		return nil, err
	}

	screp := database.Psql.ServiceConfig()
	serviceConfig, err = screp.Update(serviceConfig)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{Resp: "Update"}, nil
}

func (s *gRPCServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	screp := database.Psql.ServiceConfig()
	_, err := screp.Delete(&models.ServiceConfig{
		Service: req.ServiceName,
		Version: req.Version,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Resp: "Delete"}, nil
}
