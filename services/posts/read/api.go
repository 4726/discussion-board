package main

import (
	"fmt"
	"net"

	"github.com/4726/discussion-board/services/common"
	"github.com/4726/discussion-board/services/posts/models"
	pb "github.com/4726/discussion-board/services/posts/read/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

type Api struct {
	db   *gorm.DB
	grpc *grpc.Server
}

func NewApi(cfg Config) (*Api, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	log.Entry().Infof("connecting to database: %s", s)
	db, err := common.OpenDB("mysql", s)
	if err != nil {
		return nil, err
	}
	log.Entry().Infof("successfully connected to database: %s", s)
	// db.LogMode(true)
	db.AutoMigrate(&models.Comment{}, &models.Post{})

	opts := common.GRPCOptions{cfg.IPWhitelist, cfg.TLSCert, cfg.TLSKey}
	server, err := common.DefaultGRPCServer(opts)
	if err != nil {
		return nil, err
	}
	handlers := &Handlers{db}
	pb.RegisterPostsReadServer(server, handlers)

	return &Api{db, server}, err
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return common.RunGRPCWithGracefulShutdown(a.grpc, lis)
}
