package main

import (
	"context"
	"time"

	likes "github.com/4726/discussion-board/api-gateway/pb/likes"
	postsread "github.com/4726/discussion-board/api-gateway/pb/posts-read"
	postswrite "github.com/4726/discussion-board/api-gateway/pb/posts-write"
	search "github.com/4726/discussion-board/api-gateway/pb/search"
	user "github.com/4726/discussion-board/api-gateway/pb/user"
)

const defaultGRPCTimeout = 10

type GRPCClients struct {
	Search     search.SearchClient
	User       user.UserClient
	Likes      likes.LikesClient
	PostsRead  postsread.PostsReadClient
	PostsWrite postswrite.PostsWriteClient
}

func DefaultGRPCContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*defaultGRPCTimeout)
}
