package server

import (
	"context"
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
	return &pb.CreateResponse{Resp: "Create"}, nil
}

func (s *gRPCServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	return &pb.ReadResponse{Resp: "Read"}, nil
}

func (s *gRPCServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return &pb.UpdateResponse{Resp: "Update"}, nil
}

func (s *gRPCServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{Resp: "Delete"}, nil
}
