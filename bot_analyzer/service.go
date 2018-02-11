package bot_analyzer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/umsu2/bot_analyszer/grpc_client"
	"github.com/umsu2/bot_analyszer/grpc_service"
	"github.com/umsu2/bot_analyszer/models"
)

type gatewayService struct {
	gRPCClient grpc_client.GRPCAnalyzeService
}

type GateWayService interface {
	Analyze(context.Context, *http.Request) (grpc_service.GeneralResponse, error)
}

func NewGatewayService(gRPCClient grpc_client.GRPCAnalyzeService) gatewayService {
	return gatewayService{
		gRPCClient: gRPCClient,
	}
}

func (s gatewayService) Analyze(ctx context.Context, r *http.Request) (grpc_service.GeneralResponse, error) {
	fmt.Printf("analyze stuff: %+v \n", r)
	// get ip address etc. find routes etc.
	//bytes, err := httputil.DumpRequest(r, true)

	fmt.Printf("%+v \n", r.URL.RawQuery) // get parameters
	fmt.Printf("%+v \n", r.URL.Path)     // url path
	fmt.Printf("%+v \n", r.Method)       // url path
	fmt.Printf("%+v \n", r.Host)         // host it tries to connect
	fmt.Printf("%+v \n", r.RemoteAddr)   // addresss of requester
	fmt.Printf("%+v \n", r.Header)       // header, which is map of strings
	req := models.AppgatewayWebRequest{}

	return s.gRPCClient.Analyze(ctx, req)

}

type GateServiceMiddleware func(GateWayService) GateWayService
