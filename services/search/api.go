package main

import (
	"net"

	"github.com/4726/discussion-board/services/common"
	pb "github.com/4726/discussion-board/services/search/pb"
	"google.golang.org/grpc"
)

type Api struct {
	esc  *ESClient
	grpc *grpc.Server
}

func NewApi(cfg Config) (*Api, error) {
	log.Entry().Infof("connecting to elasticsearch: %s", cfg.ESAddr)
	esc, err := NewESClient(cfg.ESIndex, cfg.ESAddr)
	if err != nil {
		return nil, err
	}
	log.Entry().Infof("successfully connected to elasticsearch: %s", cfg.ESAddr)

	opts := common.GRPCOptions{cfg.IPWhitelist, cfg.TLSCert, cfg.TLSKey, log.Entry()}
	server, err := common.DefaultGRPCServer(opts)
	if err != nil {
		return nil, err
	}
	handlers := &Handlers{esc}
	pb.RegisterSearchServer(server, handlers)

	return &Api{esc, server}, err
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Entry().Infof("server running on addr: %s", addr)

	return common.RunGRPCWithGracefulShutdown(a.grpc, lis)
}
