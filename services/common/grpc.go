package common

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCOptions struct {
	IPWhitelist     []string
	TLSCert, TLSKey string
	LogEntry        *logrus.Entry
}

func DefaultGRPCServer(opts GRPCOptions) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(opts.TLSCert, opts.TLSKey)
	if err != nil {
		return nil, err
	}

	return grpc.NewServer(
		grpc.Creds(creds),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			IPWhiteListUnaryServerInterceptor(opts.IPWhitelist),
			otgrpc.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(opts.LogEntry),
		),
	), nil
}

func RunGRPCWithGracefulShutdown(s *grpc.Server, lis net.Listener) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	shutdownCh := make(chan error, 1)
	go func() {
		sig := <-c
		s.GracefulStop()
		shutdownCh <- fmt.Errorf(sig.String())
	}()

	serveCh := make(chan error, 1)
	go func() {
		err := s.Serve(lis)
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		return err
	case err := <-shutdownCh:
		return err
	}
}
