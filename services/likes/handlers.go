package main

import (
	"context"
	"github.com/4726/discussion-board/services/likes/pb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"time"
)

type GRPCHandlers struct {
	db *gorm.DB
}

func (h *GRPCHandlers) LikePost(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	like := PostLike{*idu.Id, *idu.UserId, time.Now()}

	if err := h.db.FirstOrCreate(&PostLike{}, &like).Error; err != nil {
		return nil, err
	}

	var count uint64

	if err := h.db.Where("post_id = ?", like.PostID).Find(&PostLike{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &pb.Total{Total: &count}, nil
}

func (h *GRPCHandlers) UnlikePost(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	like := PostLike{PostID: *idu.Id, UserID: *idu.UserId}

	if err := h.db.Delete(&like).Error; err != nil {
		return nil, err
	}

	var count uint64

	if err := h.db.Where("post_id = ?", like.PostID).Find(&PostLike{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &pb.Total{Total: &count}, nil
}

func (h *GRPCHandlers) LikeComment(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	like := CommentLike{*idu.Id, *idu.UserId, time.Now()}

	if err := h.db.FirstOrCreate(&CommentLike{}, &like).Error; err != nil {
		return nil, err
	}

	var count uint64

	if err := h.db.Where("comment_id = ?", like.CommentID).Find(&CommentLike{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &pb.Total{Total: &count}, nil
}

func (h *GRPCHandlers) UnlikeComment(ctx context.Context, idu *pb.IDUserID) (*pb.Total, error) {
	like := CommentLike{CommentID: *idu.Id, UserID: *idu.UserId}

	if err := h.db.Delete(&like).Error; err != nil {
		return nil, err
	}

	var count uint64

	if err := h.db.Where("comment_id = ?", like.CommentID).Find(&CommentLike{}).Count(&count).Error; err != nil {
		return nil, err
	}

	return &pb.Total{Total: &count}, nil
}

func (h *GRPCHandlers) GetPostLikes(ctx context.Context, ids *pb.IDs) (*pb.TotalLikes, error) {
	likes := []*pb.TotalLikes_IDLikes{}

	for _, v := range ids.Id {
		var count uint64

		if err := h.db.Where("post_id = ?", v).Find(&PostLike{}).Count(&count).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: proto.Uint64(0)})
				continue
			}
			return nil, err
		}

		likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: &count})
	}

	return &pb.TotalLikes{IdLikes: likes}, nil
}

func (h *GRPCHandlers) GetCommentLikes(ctx context.Context, ids *pb.IDs) (*pb.TotalLikes, error) {
	likes := []*pb.TotalLikes_IDLikes{}

	for _, v := range ids.Id {
		var count uint64

		if err := h.db.Where("comment_id = ?", v).Find(&CommentLike{}).Count(&count).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: proto.Uint64(0)})
				continue
			}
			return nil, err
		}

		likes = append(likes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(v), Total: &count})
	}

	return &pb.TotalLikes{IdLikes: likes}, nil
}

func (h *GRPCHandlers) PostsHaveLike(ctx context.Context, idu *pb.IDsUserID) (*pb.HaveLikes, error) {
	likes := []*pb.HaveLikes_HaveLike{}

	for _, v := range idu.Id {
		like := PostLike{}

		if err := h.db.Where("post_id = ? AND user_id = ?", v, *idu.UserId).Find(&like).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(false)})
				continue
			}
			return nil, err
		}

		likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(true)})
	}

	return &pb.HaveLikes{HaveLikes: likes}, nil
}

func (h *GRPCHandlers) CommentsHaveLike(ctx context.Context, idu *pb.IDsUserID) (*pb.HaveLikes, error) {
	likes := []*pb.HaveLikes_HaveLike{}

	for _, v := range idu.Id {
		like := CommentLike{}

		if err := h.db.Where("comment_id = ? AND user_id = ?", v, *idu.UserId).Find(&like).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(false)})
				continue
			}
			return nil, err
		}

		likes = append(likes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(v), HasLike: proto.Bool(true)})
	}

	return &pb.HaveLikes{HaveLikes: likes}, nil
}
