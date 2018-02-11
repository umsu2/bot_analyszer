package main

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"github.com/umsu2/bot_analyszer/grpc_client"
	"github.com/umsu2/bot_analyszer/grpc_service"
	"github.com/umsu2/bot_analyszer/models"
)

const (
	port = ":50051"
)

const gRPC_Server_Port = ":50051"

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", gRPC_Server_Port, "caller", log.DefaultCaller)
	var grpcServer grpc_service.WebRequestServiceServer

	ep := grpc_client.NewServerEndpoints(NewBasicService(), logger)
	grpcServer = NewGRPCServer(ep, logger)
	var g group.Group
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", gRPC_Server_Port)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", gRPC_Server_Port)
			baseServer := grpc.NewServer()
			grpc_service.RegisterWebRequestServiceServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}

	logger.Log("exit", g.Run())
}

type grpcServer struct {
	analyse grpctransport.Handler
}

type basicService struct{}

func NewBasicService() grpc_client.GRPCAnalyzeService {
	return basicService{}
}

func (s basicService) Analyze(ctx context.Context, req models.AppgatewayWebRequest) (grpc_service.GeneralResponse, error) {
	fmt.Println("FINALLY DOING FUCKING SHIT HERE, THIS IS SERVICE")
	return grpc_service.GeneralResponse{Success: true}, nil
}

func (s *grpcServer) Anaylse(ctx context.Context, req *grpc_service.WebRequest) (*grpc_service.GeneralResponse, error) {

	fmt.Println("FINALLY RECIEVED MESSAGE")
	fmt.Println(req)
	_, resp, err := s.analyse.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("+v%", resp)
	return &grpc_service.GeneralResponse{Success: true}, nil
}

func NewGRPCServer(endpoints grpc_client.EnpointSet, logger log.Logger) grpc_service.WebRequestServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		analyse: grpctransport.NewServer(
			endpoints.AnalyzeEndpoint,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
	}
}

func decodeGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*grpc_service.WebRequest)
	return models.AppgatewayWebRequest{req.RemoteIPAddress}, nil
}

func encodeGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(grpc_service.GeneralResponse)
	return &resp, nil
}
