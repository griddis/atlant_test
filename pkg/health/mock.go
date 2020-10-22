package health

import (
	"context"
)

type mockService struct {
	err error
}

func NewMockService(err error) *mockService {
	return &mockService{err}
}

func (s *mockService) GetLiveness(ctx context.Context, req *GetLivenessRequest) (resp *GetLivenessRespone, err error) {
	if s.err != nil {
		return &GetLivenessRespone{}, s.err
	} else {
		return &GetLivenessRespone{}, nil
	}
}

func (s *mockService) GetReadiness(ctx context.Context, req *GetReadinessRequest) (resp *GetReadinessResponse, err error) {
	if s.err != nil {
		return &GetReadinessResponse{}, s.err
	} else {
		return &GetReadinessResponse{}, nil
	}
}

func (s *mockService) GetVersion(ctx context.Context, req *GetVersionRequest) (resp *GetVersionResponse, err error) {
	if s.err != nil {
		return &GetVersionResponse{}, s.err
	} else {
		return &GetVersionResponse{}, nil
	}
}
