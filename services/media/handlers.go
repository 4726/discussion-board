package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/4726/discussion-board/services/media/pb"
	"github.com/golang/protobuf/proto"
	"github.com/minio/minio-go/v6"
	"github.com/segmentio/ksuid"
)

type Handlers struct {
	mc *minio.Client
}

func (h *Handlers) Upload(ctx context.Context, in *pb.UploadRequest) (*pb.Name, error) {
	buffer := bytes.NewBuffer(in.Media)
	guid, err := ksuid.NewRandom() //not guaranteed unique
	if err != nil {
		return nil, err
	}
	name := guid.String()
	opts := minio.PutObjectOptions{ContentType: "text/plain"}
	_, err = h.mc.PutObject(bucketName, name, buffer, int64(len(in.Media)), opts)
	if err != nil {
		return nil, err
	}

	return &pb.Name{Name: proto.String(name)}, nil
}

func (h *Handlers) Remove(ctx context.Context, in *pb.Name) (*pb.RemoveResponse, error) {
	if err := h.mc.RemoveObject(bucketName, in.GetName()); err != nil {
		return nil, err
	}
	return &pb.RemoveResponse{}, nil
}

func (h *Handlers) Info(ctx context.Context, in *pb.InfoRequest) (*pb.InfoResponse, error) {
	addr := fmt.Sprintf("%s/%s/", h.mc.EndpointURL().String(), bucketName)
	return &pb.InfoResponse{StoreAddress: proto.String(addr)}, nil
}
