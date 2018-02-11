package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/umsu2/bot_analyszer/bot_analyzer"
	"github.com/umsu2/bot_analyszer/grpc_client"
	"google.golang.org/grpc"
)

//go:generate protoc -I ./grpc_service --go_out=plugins=grpc:./grpc_service ./grpc_service/grpc_service.proto

const gRPC_Service_Address = "localhost:50051"

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	// todo: instrumenting code
	var sv bot_analyzer.GateWayService

	conn, err := grpc.Dial(gRPC_Service_Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Log("error", "did not connect: %v", err)
	}
	grcpClient := grpc_client.NewGRPCClient(conn, logger)
	sv = bot_analyzer.NewGatewayService(grcpClient)
	sv = bot_analyzer.LoggingMiddleware(logger)(sv)

	handler := httptransport.NewServer(
		makeEndpoint(sv),
		decodeRequest,
		encodeResponse,
	)

	http.Handle("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}

func makeEndpoint(sv bot_analyzer.GateWayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*http.Request)
		resp, err := sv.Analyze(ctx, req)
		return resp, err
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusOK)
	return nil
}
