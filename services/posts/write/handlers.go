package main

import (
	"context"
	"fmt"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/4726/discussion-board/services/posts/write/pb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var (
	PostDoesNotExistError = fmt.Errorf("post does not exist")
)

type Handlers struct {
	db *gorm.DB
}

func (h *Handlers) CreatePost(ctx context.Context, in *pb.PostRequest) (*pb.PostId, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	created := time.Now()
	post := models.Post{
		UserID:    in.GetUserId(),
		Title:     in.GetTitle(),
		Body:      in.GetBody(),
		Likes:     0,
		CreatedAt: created,
		UpdatedAt: created,
	}

	if err := h.db.Save(&post).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PostId{PostId: proto.Uint64(post.ID)}, nil
}

func (h *Handlers) DeletePost(ctx context.Context, in *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	post := models.Post{ID: in.GetPostId()}

	if in.GetUserId() != 0 {
		if err := h.db.Where("user_id = ?", in.GetUserId()).Delete(&post).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, status.Error(codes.InvalidArgument, "user did not create this post")
			}
			return nil, status.Error(codes.Internal, err.Error())
		} else {
			return &pb.DeletePostResponse{}, nil
		}
	}

	if err := h.db.Delete(&post).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "post not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *Handlers) SetPostLikes(ctx context.Context, in *pb.SetLikes) (*pb.SetLikesResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	post := models.Post{ID: in.GetId()}

	//uses UpdateColumn() instead of Update() because Update()
	//automatically updates the UpdatedAt field
	if err := h.db.Model(&post).UpdateColumn("Likes", in.GetLikes()).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "post not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SetLikesResponse{}, nil
}

func (h *Handlers) CreateComment(ctx context.Context, in *pb.CommentRequest) (*pb.CreateCommentResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	created := time.Now()
	comment := models.Comment{
		PostID:    in.GetPostId(),
		ParentID:  in.GetParentId(),
		UserID:    in.GetUserId(),
		Body:      in.GetBody(),
		CreatedAt: created,
		Likes:     0,
	}

	if err := addCommentToDB(h.db, &comment); err != nil {
		if gorm.IsRecordNotFoundError(err) || err == PostDoesNotExistError {
			return nil, status.Error(codes.NotFound, "post not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateCommentResponse{}, nil
}

func (h *Handlers) ClearComment(ctx context.Context, in *pb.ClearCommentRequest) (*pb.ClearCommentResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	comment := models.Comment{ID: in.GetCommentId()}

	if in.GetUserId() != 0 {
		if err := h.db.Model(&comment).Where("user_id = ?", in.GetUserId()).UpdateColumn("Body", "").Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, status.Error(codes.InvalidArgument, "comment does not belong to user")
			}
			return nil, status.Error(codes.Internal, err.Error())
		} else {
			return &pb.ClearCommentResponse{}, nil
		}
	}

	if err := h.db.Model(&comment).UpdateColumn("Body", "").Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "comment not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ClearCommentResponse{}, nil
}

func (h *Handlers) SetCommentLikes(ctx context.Context, in *pb.SetLikes) (*pb.SetLikesResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	comment := models.Comment{ID: in.GetId()}

	if err := h.db.Model(&comment).Update("Likes", in.GetLikes()).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "comment not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SetLikesResponse{}, nil
}

func addCommentToDB(db *gorm.DB, comment *models.Comment) error {
	post := models.Post{ID: comment.PostID}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	count := 0
	if err := tx.Model(&post).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count < 1 {
		tx.Rollback()
		return PostDoesNotExistError
	}
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&post).Update("UpdatedAt", comment.CreatedAt).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
