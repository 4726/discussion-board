package common

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func IPWhiteListUnaryServerInterceptor(whitelist []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		peer, ok := peer.FromContext(ctx)
		host, _, err := net.SplitHostPort(peer.Addr.String())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "%s is not whitelisted.", host)
		}
		var whitelisted bool
		for _, ip := range whitelist {
			if host == ip {
				whitelisted = true
				break
			}
		}
		if !whitelisted {
			return nil, status.Errorf(codes.Unauthenticated, "%s is not whitelisted.", host)
		}
		return handler(ctx, req)
	}
}
