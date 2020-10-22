package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "github.com/griddis/atlant_test/cmd/service/pb"
	"github.com/griddis/atlant_test/tools/logging"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	fetch grpctransport.Handler
	list  grpctransport.Handler
}

type ContextGRPCKey struct{}

type GRPCInfo struct{}

// NewGRPCServer makes a set of endpoints available as a gRPC Server.
func NewGRPCServer(s Service, logger log.Logger) pb.ServiceServer {
	options := []grpctransport.ServerOption{
		// grpctransport.ServerErrorLogger(logger),
		grpctransport.ServerBefore(grpcToContext()),
	}

	return &grpcServer{
		fetch: grpctransport.NewServer(
			makeFetchEndpoint(s),
			decodeGRPCFetchRequest,
			encodeGRPCFetchResponse,
			options...,
		),
		list: grpctransport.NewServer(
			makeListEndpoint(s),
			decodeGRPCListRequest,
			encodeGRPCListResponse,
			options...,
		),
	}
}

func JoinGRPC(ctx context.Context, s Service) func(*googlegrpc.Server) {
	logger := logging.FromContext(ctx)
	logger = logger.With("grpc handler", "health")
	return func(g *googlegrpc.Server) {
		pb.RegisterServiceServer(g, NewGRPCServer(s, logger))
	}
}

func grpcToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		return context.WithValue(ctx, ContextGRPCKey{}, GRPCInfo{})
	}
}

func (s *grpcServer) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	_, rep, err := s.fetch.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FetchResponse), nil
}

func (s *grpcServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	_, rep, err := s.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListResponse), nil
}

func decodeGRPCFetchRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.FetchRequest)
	if !ok {
		return nil, errors.New("decodeGRPCFetchRequest wrong request")
	}

	req := PBToFetchRequest(inReq)
	return *req, nil
}

func decodeGRPCListRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.ListRequest)
	if !ok {
		return nil, errors.New("decodeGRPCListRequest wrong request")
	}

	req := PBToListRequest(inReq)
	return *req, nil
}

func encodeGRPCFetchResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*FetchResponse)
	if !ok {
		return nil, errors.New("encodeGRPCFetchResponse wrong response")
	}

	return FetchResponseToPB(inResp), nil
}

func encodeGRPCListResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*ListResponse)
	if !ok {
		return nil, errors.New("encodeGRPCListResponse wrong response")
	}

	return ListResponseToPB(inResp), nil
}

func FetchRequestToPB(d *FetchRequest) *pb.FetchRequest {
	if d == nil {
		return nil
	}

	resp := pb.FetchRequest{}

	for _, c := range d.Files {
		v := c
		resp.Files = append(resp.Files, v)
	}

	return &resp
}

func PBToFetchRequest(d *pb.FetchRequest) *FetchRequest {
	if d == nil {
		return nil
	}

	resp := FetchRequest{}

	for _, c := range d.Files {
		v := c
		resp.Files = append(resp.Files, v)
	}

	return &resp
}

func FetchResponseToPB(d *FetchResponse) *pb.FetchResponse {
	if d == nil {
		return nil
	}

	resp := pb.FetchResponse{
		Status: d.Status,
	}

	return &resp
}

func PBToFetchResponse(d *pb.FetchResponse) *FetchResponse {
	if d == nil {
		return nil
	}

	resp := FetchResponse{
		Status: d.Status,
	}

	return &resp
}

func ListItemsToPB(d *ListItems) *pb.ListItems {
	if d == nil {
		return nil
	}

	resp := pb.ListItems{
		Id:      d.Id,
		Price:   d.Price,
		Counter: d.Counter,
		Date:    d.Date.String(),
	}

	return &resp
}

func PBToListItems(d *pb.ListItems) *ListItems {
	if d == nil {
		return nil
	}

	t, _ := time.Parse(d.Date, "2020-10-21T19:01:25.962+00:00")
	resp := ListItems{
		Id:      d.Id,
		Price:   d.Price,
		Counter: d.Counter,
		Date:    t,
	}

	return &resp
}

func ListRequestToPB(d *ListRequest) *pb.ListRequest {
	if d == nil {
		return nil
	}

	resp := pb.ListRequest{
		Limiter: LimiterToPB(d.Limiter),
	}

	for i, c := range d.Sorter {
		resp.Sorter[i] = c
	}

	return &resp
}

func PBToListRequest(d *pb.ListRequest) *ListRequest {
	if d == nil {
		return nil
	}

	resp := ListRequest{
		Limiter: PBToLimiter(d.Limiter),
	}

	for i, c := range d.Sorter {
		resp.Sorter[i] = c
	}

	return &resp
}

func LimiterToPB(d *Limiter) *pb.Limiter {
	if d == nil {
		return nil
	}

	resp := pb.Limiter{
		Offsetbyid: d.Offsetbyid,
		Limit:      d.Limit,
	}

	return &resp
}

func PBToLimiter(d *pb.Limiter) *Limiter {
	if d == nil {
		return nil
	}

	resp := Limiter{
		Offsetbyid: d.Offsetbyid,
		Limit:      d.Limit,
	}

	return &resp
}

func ListResponseToPB(d *ListResponse) *pb.ListResponse {
	if d == nil {
		return nil
	}

	resp := pb.ListResponse{}

	for _, c := range d.ListItems {
		v := c
		resp.ListItems = append(resp.ListItems, ListItemsToPB(v))
	}

	return &resp
}

func PBToListResponse(d *pb.ListResponse) *ListResponse {
	if d == nil {
		return nil
	}

	resp := ListResponse{}

	for _, c := range d.ListItems {
		v := c
		resp.ListItems = append(resp.ListItems, PBToListItems(v))
	}

	return &resp
}
