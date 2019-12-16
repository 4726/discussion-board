package common

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCOptions struct {
	IPWhitelist     []string
	TLSCert, TLSKey string
}

func DefaultGRPCServer(opts GRPCOptions) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(opts.TLSCert, opts.TLSKey)
	if err != nil {
		return nil, err
	}

	return grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			IPWhiteListUnaryServerInterceptor(opts.IPWhitelist),
			otgrpc.UnaryServerInterceptor(),
		)),
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
