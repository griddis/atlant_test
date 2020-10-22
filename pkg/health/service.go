package health

import (
	"context"
	"errors"
)

var (
	BuildTime          string
	Commit             string
	Version            string
	ErrServiceNotReady = errors.New("service not ready")
)

type canBeReady interface {
	IsReady() bool
}

func NewService(services ...canBeReady) Service {
	return &service{services}
}

type service struct {
	// app is ready if cache loaded
	services []canBeReady
}

func (s *service) GetLiveness(ctx context.Context, req *GetLivenessRequest) (*GetLivenessRespone, error) {
	return &GetLivenessRespone{Status: "ok"}, nil
}

func (s *service) GetReadiness(ctx context.Context, req *GetReadinessRequest) (*GetReadinessResponse, error) {
	for _, r := range s.services {
		if !r.IsReady() {
			return nil, ErrServiceNotReady
		}
	}
	return &GetReadinessResponse{Status: "ok"}, nil
}

func (s *service) GetVersion(ctx context.Context, req *GetVersionRequest) (*GetVersionResponse, error) {
	return &GetVersionResponse{Version: RespVersion{
		Commit,
		Version,
		BuildTime,
	}}, nil
}
