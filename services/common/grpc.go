package common

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
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
