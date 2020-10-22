package integration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/griddis/atlant_test/cmd/service"
	"github.com/griddis/atlant_test/configs"
	"github.com/griddis/atlant_test/internal/server"
	"github.com/griddis/atlant_test/pkg/repository"
	"github.com/griddis/atlant_test/tools/logging"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

const grpcAddrService = "localhost:10081"

var (
	onceInit sync.Once
	client   = NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		//assert.Equal(t, req.URL.String(), "https://gist.githubusercontent.com/griddis/cf62cdfaa46d779dd6f7f7b436ceb77d/raw/e143b62d01bef89d77f3a9e1e1001d73027da224/gistfile1.txt")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`test1,0.31\ntest2,4.32\ntest3,3.11`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
)

func RunTestServer(client *http.Client) {
	onceInit.Do(func() {
		cfg := configs.NewConfig()
		if err := cfg.Read(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read config: %s", err)
			os.Exit(1)
		}
		if err := cfg.Print(); err != nil {
			fmt.Fprintf(os.Stderr, "print config: %s", err)
			os.Exit(1)
		}

		var (
			logger = logging.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx = logging.WithContext(ctx, logger)

		rctx, _ := context.WithTimeout(ctx, 10*time.Second)
		cfg.Database.Driver = "inmemory"
		cfg.Server.GRPC.Port = 10081
		repository, err := repository.NewRepo(rctx, logger, &cfg.Database)

		if err != nil {
			logger.Error("init", "repository", "err", err)
			os.Exit(1)
		}
		defer repository.Close(rctx)

		mainService := service.NewService(ctx, cfg, repository, client)

		s, err := server.NewServer(
			server.SetConfig(&cfg.Server),
			server.SetLogger(logger),
			server.SetGRPC(
				service.JoinGRPC(ctx, mainService),
			),
		)
		if err != nil {
			logger.Error("init", "server", "err", err)
			os.Exit(1)
		}
		//defer s.Close()

		if err = s.AddGRPC(); err != nil {
			logger.Error("err", err)
			os.Exit(1)
		}

		s.AddSignalHandler()
		go func() {
			s.Run()
		}()
	})

}
func TestGRPCServiceFetchWithFile(t *testing.T) {

	RunTestServer(client)

	conn, err := grpc.Dial(grpcAddrService, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connection to grpc server: %s", err)
	}
	defer conn.Close()

	grpcclient := service.NewGRPCClient(conn, log.NewNopLogger())
	_, err = grpcclient.Fetch(context.Background(), &service.FetchRequest{Files: []string{"https://gist.githubusercontent.com/griddis/cf62cdfaa46d779dd6f7f7b436ceb77d/raw/e143b62d01bef89d77f3a9e1e1001d73027da224/gistfile1.txt"}})

	assert.NoError(t, err)
}

func TestGRPCServiceFetchEmpty(t *testing.T) {

	RunTestServer(client)

	conn, err := grpc.Dial(grpcAddrService, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connection to grpc server: %s", err)
	}
	defer conn.Close()

	grpcclient := service.NewGRPCClient(conn, log.NewNopLogger())
	_, err = grpcclient.Fetch(context.Background(), &service.FetchRequest{})

	assert.NoError(t, err)
}
