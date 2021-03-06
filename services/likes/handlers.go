package main

import (
	"context"
	"time"

	pb "github.com/4726/discussion-board/services/likes/pb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handlers struct {
	db *gorm.DB
}

func (h *Handlers) LikePost(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	like := PostLike{idu.GetId(), idu.GetUserId(), time.Now()}

	if err := h.db.FirstOrCreate(&PostLike{}, &like).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var count uint64

	if err := h.db.Model(&PostLike{}).Where("post_id = ?", like.PostID).Count(&count).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Total{Total: proto.Uint64(count)}, nil
}

func (h *Handlers) UnlikePost(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	like := PostLike{PostID: idu.GetId(), UserID: idu.GetUserId()}

	if err := h.db.Delete(&like).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var count uint64

	if err := h.db.Model(&PostLike{}).Where("post_id = ?", like.PostID).Count(&count).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.Total{Total: proto.Uint64(count)}, nil
}

func (h *Handlers) LikeComment(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	like := CommentLike{idu.GetId(), idu.GetUserId(), time.Now()}

	if err := h.db.FirstOrCreate(&CommentLike{}, &like).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var count uint64

	if err := h.db.Model(&CommentLike{}).Where("comment_id = ?", like.CommentID).Count(&count).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Total{Total: proto.Uint64(count)}, nil
}

func (h *Handlers) UnlikeComment(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	like := CommentLike{CommentID: idu.GetId(), UserID: idu.GetUserId()}

	if err := h.db.Delete(&like).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var count uint64

	if err := h.db.Model(&CommentLike{}).Where("comment_id = ?", like.CommentID).Count(&count).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.Total{Total: proto.Uint64(count)}, nil
}

func (h *Handlers) GetPostLikes(ctx context.Context, ids *pb.IDs) (*pb.TotalLikes, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	likes := []*pb.TotalLikes_IDLikes{}

	for _, v := range ids.Id {
		var count uint64

		if err := h.db.Model(&PostLike{}).Where("post_id = ?", v).Count(&count).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: proto.Uint64(0)})
				continue
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: &count})
	}

	return &pb.TotalLikes{IdLikes: likes}, nil
}

func (h *Handlers) GetCommentLikes(ctx context.Context, ids *pb.IDs) (*pb.TotalLikes, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	likes := []*pb.TotalLikes_IDLikes{}

	for _, v := range ids.Id {
		var count uint64

		if err := h.db.Model(&CommentLike{}).Where("comment_id = ?", v).Count(&count).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: proto.Uint64(0)})
				continue
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: &count})
	}

	return &pb.TotalLikes{IdLikes: likes}, nil
}

func (h *Handlers) PostsHaveLike(ctx context.Context, idu *pb.IDsUserID) (*pb.HaveLikes, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	likes := []*pb.HaveLikes_HaveLike{}

	for _, v := range idu.Id {
		like := PostLike{}

		if err := h.db.Where("post_id = ? AND user_id = ?", v, idu.GetUserId()).Find(&like).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(false)})
				continue
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(true)})
	}

	return &pb.HaveLikes{HaveLikes: likes}, nil
}

func (h *Handlers) CommentsHaveLike(ctx context.Context, idu *pb.IDsUserID) (*pb.HaveLikes, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	likes := []*pb.HaveLikes_HaveLike{}

	for _, v := range idu.Id {
		like := CommentLike{}

		if err := h.db.Where("comment_id = ? AND user_id = ?", v, idu.GetUserId()).Find(&like).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(false)})
				continue
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(true)})
	}

	return &pb.HaveLikes{HaveLikes: likes}, nil
}

func (h *Handlers) DeletePost(ctx context.Context, in *pb.Id) (*pb.DeletePostResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	if err := h.db.Exec("DELETE from post_likes WHERE post_id = ?", in.GetId()).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *Handlers) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING.Enum()}, nil
}
