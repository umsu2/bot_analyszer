package bot_analyzer

import (
	"context"
	"fmt"
	"net/http"
)

type gatewayService struct {
}

type GateWayService interface {
	Analyze(context.Context, *http.Request) error
}

func NewGatewayService() gatewayService {
	return gatewayService{}
}

func (s gatewayService) Analyze(ctx context.Context, r *http.Request) error {
	fmt.Printf("analyze stuff: %+v \n", r)
	// get ip address etc. find routes etc.
	//bytes, err := httputil.DumpRequest(r, true)

	fmt.Printf("%+v \n", r.URL.RawQuery) // get parameters
	fmt.Printf("%+v \n", r.URL.Path)     // url path
	fmt.Printf("%+v \n", r.Method)       // url path
	fmt.Printf("%+v \n", r.Host)         // host it tries to connect
	fmt.Printf("%+v \n", r.RemoteAddr)   // addresss of requester
	fmt.Printf("%+v \n", r.Header)       // header, which is map of strings

	// pass the request data, and body into the rpc methods
	return s.next.Analyze(ctx, r)
}

type AppgatewayWebRequest struct {
	IPAddress string
}

type GateServiceMiddleware func(GateWayService) GateWayService
