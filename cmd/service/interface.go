package service

import (
	"context"
)

type Service interface {
	Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error)
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)
}
