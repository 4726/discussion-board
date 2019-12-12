package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/4726/discussion-board/services/search/pb"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

type Api struct {
	esc  *ESClient
	grpc *grpc.Server
}

func NewApi(cfg Config) (*Api, error) {
	esc, err := NewESClient(cfg.ESIndex, cfg.ESAddr)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.UnaryServerInterceptor()))
	handlers := &Handlers{esc}
	pb.RegisterSearchServer(server, handlers)

	return &Api{esc, server}, err
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	shutdownCh := make(chan error, 1)
	go func() {
		sig := <-c
		a.grpc.GracefulStop()
		shutdownCh <- fmt.Errorf(sig.String())
	}()

	serveCh := make(chan error, 1)
	go func() {
		err := a.grpc.Serve(lis)
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		return err
	case err := <-shutdownCh:
		return err
	}
}
