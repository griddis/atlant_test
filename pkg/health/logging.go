package health

import (
	"context"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/griddis/atlant_test/tools/logging"
)

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(ctx context.Context, s Service) Service {
	logger := logging.FromContext(ctx)
	logger = logger.With("component", "health")
	return &loggingService{logger, s}
}

type loggingService struct {
	logger *logging.Logger
	Service
}

func (s *loggingService) GetLiveness(ctx context.Context, req *GetLivenessRequest) (resp *GetLivenessRespone, err error) {
	defer func(begin time.Time) {
		m := getInfoFromContext(ctx)
		m = append(m,
			"code", getHTTPStatusCode(err),
			"method", "GetLiveness",
			"took", time.Since(begin),
		)

		if getHTTPStatusCode(err) == 404 {
			m = append(m, "msg", err)
			level.Warn(s.logger).Log(m...)
		} else if err != nil {
			m = append(m, "err", err)
			level.Error(s.logger).Log(m...)
		} else {
			level.Info(s.logger).Log(m...)
		}
	}(time.Now())
	return s.Service.GetLiveness(ctx, req)
}

func (s *loggingService) GetReadiness(ctx context.Context, req *GetReadinessRequest) (resp *GetReadinessResponse, err error) {
	defer func(begin time.Time) {
		m := getInfoFromContext(ctx)
		m = append(m,
			"code", getHTTPStatusCode(err),
			"method", "GetReadiness",
			"took", time.Since(begin),
		)

		if getHTTPStatusCode(err) == 404 {
			m = append(m, "msg", err)
			level.Warn(s.logger).Log(m...)
		} else if err != nil {
			m = append(m, "err", err)
			level.Error(s.logger).Log(m...)
		} else {
			level.Info(s.logger).Log(m...)
		}
	}(time.Now())
	return s.Service.GetReadiness(ctx, req)
}

func (s *loggingService) GetVersion(ctx context.Context, req *GetVersionRequest) (resp *GetVersionResponse, err error) {
	defer func(begin time.Time) {
		m := getInfoFromContext(ctx)
		m = append(m,
			"code", getHTTPStatusCode(err),
			"method", "GetVersion",
			"took", time.Since(begin),
		)

		if getHTTPStatusCode(err) == 404 {
			m = append(m, "msg", err)
			level.Warn(s.logger).Log(m...)
		} else if err != nil {
			m = append(m, "err", err)
			level.Error(s.logger).Log(m...)
		} else {
			level.Info(s.logger).Log(m...)
		}
	}(time.Now())
	return s.Service.GetVersion(ctx, req)
}

func getInfoFromContext(ctx context.Context) []interface{} {
	m := make([]interface{}, 0)

	{
		val := ctx.Value(ContextGRPCKey{})
		if _, ok := val.(GRPCInfo); ok {
			m = append(m, "protocol", "GRPC")
		}
	}

	{
		val := ctx.Value(ContextHTTPKey{})
		if i, ok := val.(HTTPInfo); ok {
			m = append(m,
				// "protocol", i.Protocol,
				// "http_method", i.Method,
				// "from", i.From,
				"url", i.URL,
			)
		}
	}

	return m
}
