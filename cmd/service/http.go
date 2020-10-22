package service

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/griddis/atlant_test/tools/logging"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidArgument is returned when one or more arguments are invalid.
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	errBadRoute        = errors.New("bad route")
)

type ContextHTTPKey struct{}

type HTTPInfo struct {
	Method   string
	URL      string
	From     string
	Protocol string
}

func MakeHTTPHandler(ctx context.Context, s Service) http.Handler {
	logger := logging.FromContext(ctx)
	logger = logger.With("http handler", "additional")

	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(httpToContext()),
	}

	r.Methods("POST").Path("/_fetch").Handler(httptransport.NewServer(
		makeFetchEndpoint(s),
		decodePOSTFetchRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/_list").Handler(httptransport.NewServer(
		makeListEndpoint(s),
		decodeGETListRequest,
		encodeResponse,
		options...,
	))

	return r
}

func httpToContext() httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		return context.WithValue(ctx, ContextHTTPKey{}, HTTPInfo{
			Method:   req.Method,
			URL:      req.RequestURI,
			From:     req.RemoteAddr,
			Protocol: req.Proto,
		})
	}
}

func decodePOSTFetchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request FetchRequest

	{
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, errors.Wrap(ErrInvalidArgument, err.Error())
		}
	}

	return request, nil
}

func decodeGETListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ListRequest

	{
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, errors.Wrap(ErrInvalidArgument, err.Error())
		}
	}

	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encodeError handles error from business-layer.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("X-Esp-Error", err.Error())
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")

	w.WriteHeader(getHTTPStatusCode(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

type errorCode interface {
	Code() int
}

// getHTTPStatusCode returns http status code from error.
func getHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if e, ok := err.(errorCode); ok && e.Code() != 0 {
		return e.Code()
	}

	switch errors.Cause(err) {
	case ErrInvalidArgument:
		return http.StatusBadRequest
	case ErrAlreadyExists:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
