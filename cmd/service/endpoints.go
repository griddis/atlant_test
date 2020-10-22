package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type endpoints struct {
	FetchEndpoint endpoint.Endpoint
	ListEndpoint  endpoint.Endpoint
}

func (e endpoints) Fetch(ctx context.Context, req *FetchRequest) (resp *FetchResponse, err error) {
	response, err := e.FetchEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(FetchResponse)
	return &r, err
}

func (e endpoints) List(ctx context.Context, req *ListRequest) (resp *ListResponse, err error) {
	response, err := e.ListEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := response.(ListResponse)
	return &r, err
}

func makeFetchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FetchRequest)
		return s.Fetch(ctx, &req)
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		return s.List(ctx, &req)
	}
}
