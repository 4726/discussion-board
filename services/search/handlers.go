package main

import (
	"context"
	"github.com/4726/discussion-board/services/search/pb"
)

type Handlers struct {
	esc *ESClient
}

func (h *Handlers) Index(ctx context.Context, in *pb.Post) (*pb.IndexResponse, error) {
	if ctx.Err() == context.Canceled {return nil, fmt.Errorf("client cancelled")}
	post := Post{
		in.GetTitle(),
		in.GetBody(),
		in.GetId(),
		in.GetUserId(),
		in.GetTimestamp(),
		in.GetLikes(),
	}

	if err := h.esc.Index(post); err != nil {
		return nil, err
	}
	return &pb.IndexResponse{}, nil
}

func (h *Handlers) Search(ctx context.Context, in *pb.SearchQuery) (*pb.SearchResult, error) {
	if ctx.Err() == context.Canceled {return nil, fmt.Errorf("client cancelled")}
	res, err := h.esc.Search(in.GetTerm(), in.GetFrom(), in.GetTotal())
	if err != nil {
		return nil, err
	}
	return &pb.SearchResult{Id: res}, nil
}

func (h *Handlers) SetLikes(ctx context.Context, in *pb.Likes) (*pb.LikesResponse, error) {
	if ctx.Err() == context.Canceled {return nil, fmt.Errorf("client cancelled")}
	if err := h.esc.UpdateLikes(in.GetId(), in.GetLikes()); err != nil {
		return nil, err
	}
	return &pb.LikesResponse{}, nil
}

func (h *Handlers) DeletePost(ctx context.Context, in *pb.Id) (*pb.DeletePostResponse, error) {
	if ctx.Err() == context.Canceled {return nil, fmt.Errorf("client cancelled")}
	if err := h.esc.Delete(in.GetId()); err != nil {
		return nil, err
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *Handlers) SetTimestamp(ctx context.Context, in *pb.Timestamp) (*pb.SetTimestampResponse, error) {
	if ctx.Err() == context.Canceled {return nil, fmt.Errorf("client cancelled")}
	if err := h.esc.UpdateLastUpdate(in.GetId(), in.GetTimestamp()); err != nil {
		return nil, err
	}
	return &pb.SetTimestampResponse{}, nil
}
