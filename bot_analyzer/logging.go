package bot_analyzer

import (
	"time"

	"context"

	"net/http"

	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/umsu2/bot_analyszer/grpc_service"
)

type logmw struct {
	logger log.Logger
	GateWayService
}

func LoggingMiddleware(logger log.Logger) GateServiceMiddleware {
	return func(next GateWayService) GateWayService {
		return logmw{logger, next}
	}
}

func (mw logmw) Analyze(ctx context.Context, s *http.Request) (grpc_service.GeneralResponse, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "analyze",
			"input", fmt.Sprintf("%+v", s),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.GateWayService.Analyze(ctx, s)

}
