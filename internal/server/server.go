package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/gorilla/mux"
	"github.com/griddis/atlant_test/configs"
	"github.com/griddis/atlant_test/tools/limiting"
	"github.com/griddis/atlant_test/tools/logging"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Server main struct for srms service
type Server struct {
	cfg     *configs.Server
	logger  *logging.Logger
	handler http.Handler
	grpc    *grpcServer
	group   run.Group
}

type grpcServer struct {
	listener net.Listener
	server   *grpc.Server
}

type Option func(*Server)

func (s *Server) setGroup(group run.Group) {
	s.group = group
}

// Run запускает сервер
func (s *Server) Run() error {
	return s.logger.Log("exit", s.group.Run())
}

// Close closes everything that is open and should be closed after server shutdown
func (s *Server) Close() {
	if s.grpc != nil {
		s.logger.Info("component", "GRPC server", "msg", "close connection")
		s.grpc.listener.Close()
	}
}

// AddHTTP  http server start when Server.Run()
func (s *Server) AddHTTP() error {
	addr := fmt.Sprintf(":%d", s.cfg.HTTP.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cann't add HTTP transport")
	}

	s.group.Add(func() error {
		s.logger.Info("component", "HTTP server", "addr", addr, "msg", "listening...")
		if s.cfg.HTTP.Limiter.Enabled {
			limiting.SetLimiter(limiting.NewLimiter(s.cfg.HTTP.Limiter.Limit, 1))
			s.handler = limiting.Middleware(s.handler)
		}

		httpServer := &http.Server{
			Handler:      s.handler,
			WriteTimeout: time.Second * time.Duration(s.cfg.HTTP.TimeoutSec),
		}

		return httpServer.Serve(listener)
	}, func(error) {
		listener.Close()
	})
	return nil
}

// AddGRPC  grpc server start when Server.Run()
func (s *Server) AddGRPC() error {
	addr := fmt.Sprintf(":%d", s.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	s.grpc.listener = listener
	if err != nil {
		return errors.Wrap(err, "cann't add GRPC transport")
	}

	s.group.Add(func() error {
		s.logger.Info("component", "GRPC server", "addr", addr, "msg", "listening...")
		return s.grpc.server.Serve(listener)
	}, func(error) {
		listener.Close()
	})
	return nil
}

// AddSignalHandler add listener os signal when Server.Run()
func (s *Server) AddSignalHandler() {
	ch := make(chan struct{})
	s.group.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return errors.Errorf("received signal %s", sig)
		case <-ch:
			return nil
		}
	}, func(error) {
		close(ch)
	})
}

// NewServer инициализирует сервер.
func NewServer(ops ...Option) (svc *Server, err error) {
	svc = new(Server)

	for _, o := range ops {
		o(svc)
	}

	return svc, nil
}

func SetLogger(logger *logging.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func SetConfig(cfg *configs.Server) Option {
	return func(s *Server) {
		s.cfg = cfg
	}
}

func SetHandler(handlers map[string]http.Handler) Option {

	sts := make([]string, 0)
	for s, _ := range handlers {
		sts = append(sts, s)
	}
	sort.Strings(sts)

	return func(s *Server) {
		mux := mux.NewRouter().StrictSlash(false)

		for i := len(sts) - 1; i >= 0; i-- {
			name := sts[i]

			if handler, ok := handlers[name]; ok {
				mux.PathPrefix("/" + name).Handler(handler)
			}
		}

		s.handler = mux
	}
}

func SetGRPC(joins ...func(grpc *grpc.Server)) Option {
	return func(s *Server) {
		gServer := grpc.NewServer(
			grpc.UnaryInterceptor(grpctransport.Interceptor),
			grpc.ConnectionTimeout(time.Second*time.Duration(s.cfg.GRPC.TimeoutSec)),
		)
		for _, j := range joins {
			j(gServer)
		}
		s.grpc = &grpcServer{server: gServer}
	}
}
