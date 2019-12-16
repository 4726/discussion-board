package main

import (
	"fmt"
	"net"

	"os"
	"os/signal"
	"syscall"

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

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
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
