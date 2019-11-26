package main

import (
	"context"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/4726/discussion-board/services/posts/read/pb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
)

type Handlers struct {
	db *gorm.DB
}

func (h *Handlers) GetFullPost(ctx context.Context, in *pb.Id) (*pb.Post, error) {
	var post models.Post
	if err := h.db.First(&post, in.GetId()).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &pb.Post{}, nil
		}
		return nil, err
	}

	return modelPostToProtoPost(post), nil
}

func (h *Handlers) GetPosts(ctx context.Context, in *pb.GetPostsQuery) (*pb.MultiplePosts, error) {
	var sortType string

	switch in.GetSort() {
	case "likes_desc":
		sortType = "likes desc"
	case "created_at_desc":
		sortType = "created_at desc"
	case "created_at":
		sortType = "created_at"
	case "updated_at":
		sortType = "updated_at"
	default:
		sortType = "updated_at desc"
	}

	if in.GetUserId() != 0 {
		posts, err := getPostsUser(h.db, in.GetFrom(), in.GetTotal(), in.GetUserId(), sortType)
		if err != nil {
			return nil, err
		}

		protoPosts := []*pb.Post{}
		for _, v := range posts {
			protoPosts = append(protoPosts, modelPostToProtoPost(v))
		}
		return &pb.MultiplePosts{Posts: protoPosts}, nil
	}

	posts, err := getPosts(h.db, in.GetFrom(), in.GetTotal(), sortType)
	if err != nil {
		return nil, err
	}

	protoPosts := []*pb.Post{}
	for _, v := range posts {
		protoPosts = append(protoPosts, modelPostToProtoPost(v))
	}
	return &pb.MultiplePosts{Posts: protoPosts}, nil
}

func (h *Handlers) GetPostsById(ctx context.Context, in *pb.Ids) (*pb.MultiplePosts, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user_id", "title", "likes", "created_at", "updated_at"}
	if err := h.db.Preload("Comments").Select(selectFields).
		Where(in.Id).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	protoPosts := []*pb.Post{}
	for _, v := range posts {
		protoPosts = append(protoPosts, modelPostToProtoPost(v))
	}

	return &pb.MultiplePosts{Posts: protoPosts}, nil
}

func modelPostToProtoPost(post models.Post) *pb.Post {
	protoComments := []*pb.Comment{}

	for _, v := range post.Comments {
		protoComment := modelCommentToProtoComment(v)
		protoComments = append(protoComments, protoComment)
	}

	return &pb.Post{
		Id:        proto.Uint64(post.ID),
		UserId:    proto.Uint64(post.UserID),
		Title:     proto.String(post.Title),
		Body:      proto.String(post.Body),
		Likes:     proto.Int64(post.Likes),
		CreatedAt: proto.Int64(post.CreatedAt.Unix()),
		UpdatedAt: proto.Int64(post.UpdatedAt.Unix()),
		Comments:  protoComments,
	}
}

func modelCommentToProtoComment(comment models.Comment) *pb.Comment {
	return &pb.Comment{
		Id:        proto.Uint64(comment.ID),
		PostId:    proto.Uint64(comment.PostID),
		ParentId:  proto.Uint64(comment.ParentID),
		UserId:    proto.Uint64(comment.UserID),
		Body:      proto.String(comment.Body),
		CreatedAt: proto.Int64(comment.CreatedAt.Unix()),
		Likes:     proto.Int64(comment.Likes),
	}
}

func getPosts(db *gorm.DB, from, total uint64, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user_id", "title", "likes", "created_at", "updated_at"}
	if err := db.Preload("Comments").Select(selectFields).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func getPostsUser(db *gorm.DB, from, total uint64, userID uint64, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user_id", "title", "likes", "created_at", "updated_at"}
	if err := db.Preload("Comments").Select(selectFields).
		Where("user_id = ?", userID).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
