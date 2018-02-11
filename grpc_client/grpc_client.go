package grpc_client

import (
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"context"

	"fmt"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/umsu2/bot_analyszer/grpc_service"
	"github.com/umsu2/bot_analyszer/models"
)

type EnpointSet struct {
	AnalyzeEndpoint endpoint.Endpoint
}
type GRPCAnalyzeService interface {
	Analyze(ctx context.Context, req models.AppgatewayWebRequest) (grpc_service.GeneralResponse, error)
}

func NewServerEndpoints(service GRPCAnalyzeService, logger log.Logger) EnpointSet {

	var AnalyzeEndpoint endpoint.Endpoint

	AnalyzeEndpoint = MakeAnalyzeEndpoint(service)

	return EnpointSet{
		AnalyzeEndpoint,
	}
}

func MakeAnalyzeEndpoint(s GRPCAnalyzeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.AppgatewayWebRequest)
		resp, err := s.Analyze(ctx, req)
		return resp, err
	}
}

func (s EnpointSet) Analyze(ctx context.Context, req models.AppgatewayWebRequest) (grpc_service.GeneralResponse, error) {
	resp, err := s.AnalyzeEndpoint(ctx, req)
	if err != nil {
		return grpc_service.GeneralResponse{Success: false}, err
	}
	data, success := resp.(*grpc_service.GeneralResponse)
	if !success {
		fmt.Print(data)
	}

	return *data, nil
}

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) GRPCAnalyzeService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var anyalzeEndpoint endpoint.Endpoint
	{
		anyalzeEndpoint = grpctransport.NewClient(
			conn,
			"grpc_service.WebRequestService",
			"Anaylse",
			encodeRequest,
			decodeResponse,
			grpc_service.GeneralResponse{},
		).Endpoint()
	}
	return EnpointSet{anyalzeEndpoint}

}
func encodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	// on the client take some fucking object and serialize it to some fucking webrequest object to be sent through the wire
	req := request.(models.AppgatewayWebRequest) // get a request object and cast
	return &grpc_service.WebRequest{RemoteIPAddress: req.IPAddress}, nil
}

func decodeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*grpc_service.GeneralResponse)

	return reply, nil
}
