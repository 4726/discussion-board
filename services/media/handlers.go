package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/4726/discussion-board/services/media/pb"
	"github.com/golang/protobuf/proto"
	"github.com/minio/minio-go/v6"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handlers struct {
	mc *minio.Client
}

func (h *Handlers) Upload(ctx context.Context, in *pb.UploadRequest) (*pb.Name, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	buffer := bytes.NewBuffer(in.Media)
	guid, err := ksuid.NewRandom() //not guaranteed unique
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	name := guid.String()
	opts := minio.PutObjectOptions{ContentType: "text/plain"}
	_, err = h.mc.PutObject(bucketName, name, buffer, int64(len(in.Media)), opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Name{Name: proto.String(name)}, nil
}

func (h *Handlers) Remove(ctx context.Context, in *pb.Name) (*pb.RemoveResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	if err := h.mc.RemoveObject(bucketName, in.GetName()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RemoveResponse{}, nil
}

func (h *Handlers) Info(ctx context.Context, in *pb.InfoRequest) (*pb.InfoResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	addr := fmt.Sprintf("%s/%s/", h.mc.EndpointURL().String(), bucketName)
	return &pb.InfoResponse{StoreAddress: proto.String(addr)}, nil
}

func (h *Handlers) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING.Enum()}, nil
}
