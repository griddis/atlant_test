package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "github.com/griddis/atlant_test/cmd/service/pb"
	"google.golang.org/grpc"
)

// NewGRPCClient returns an Service backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) Service {
	// global client middlewares
	options := []grpctransport.ClientOption{
		//grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	}

	return endpoints{
		// Each individual endpoint is an grpc/transport.Client (which implements
		// endpoint.Endpoint) that gets wrapped with various middlewares. If you
		// made your own client library, you'd do this work there, so your server
		// could rely on a consistent set of client behavior.
		FetchEndpoint: grpctransport.NewClient(
			conn,
			"service.Service",
			"Fetch",
			encodeGRPCFetchRequest,
			decodeGRPCFetchResponse,
			pb.FetchResponse{},
			options...,
		).Endpoint(),
		ListEndpoint: grpctransport.NewClient(
			conn,
			"service.Service",
			"List",
			encodeGRPCListRequest,
			decodeGRPCListResponse,
			pb.ListResponse{},
			options...,
		).Endpoint(),
	}
}

func encodeGRPCFetchRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*FetchRequest)
	if !ok {
		return nil, errors.New("encodeGRPCFetchRequest wrong request")
	}

	return FetchRequestToPB(inReq), nil
}

func encodeGRPCListRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("encodeGRPCListRequest wrong request")
	}

	return ListRequestToPB(inReq), nil
}

func decodeGRPCFetchResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.FetchResponse)
	if !ok {
		return nil, errors.New("decodeGRPCFetchResponse wrong response")
	}

	resp := PBToFetchResponse(inResp)

	return *resp, nil
}

func decodeGRPCListResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.ListResponse)
	if !ok {
		return nil, errors.New("decodeGRPCListResponse wrong response")
	}

	resp := PBToListResponse(inResp)

	return *resp, nil
}
