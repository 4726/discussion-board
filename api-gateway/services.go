package main

type MediaStoreInfo struct {
	StoreAddress string
}

type MediaService interface {
	Upload(file string) (name string, error)
	Remove(name string) error
	Info() (MediaStoreInfo, error)
}

type Post struct {
	ID        uint
	User      string
	Title     string
	Body      string
	Likes     int
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}

type GetManyOptions struct {
	User, SortType string
}

type CreatePostsData struct {
	Title, Body, User string
}

type CommentData struct {
	PostID, ParentID, UserID int
	Body       string
}

type PostsService interface {
	Get(postID int) (Post, error)
	GetMany(total, from int, opts GetManyOptions) ([]Post, error)
	Create(data CreatePostsData) (postID int, error)
	Delete(postID int) error
	UpdateLikes(postID, likes int) error
	CreateComment(data CommentData) error
	ClearComment(commentID int) error
	UpdateCommentLikes(commentID, likes int) error
	DeleteIfOwner(postID, userID int) error //todo: can just accept an optional userid field
	LikePost(postID, userID int) error  //todo: userID because might already liked
	UnlikePost(postID, userID int) error //todo
	LikeComment(commentID, userID int) //todo
	UnlikeComment(commentID, userID int) //todo
	ClearCommentIfOwner(commentID int) error //todo: same as deleteifowner
	GetMultiple(postIDs []int) ([]Post, error) //todo
}

type SearchService interface {
	Search(from, total int, term string) (postIDs []int, error)
}

type Profile struct {
	UserID   int    
	Username string
	Bio      string
	AvatarID string
}

type UpdateProfileOptions struct {
	Bio, AvatarID string
}

type UserService interface {
	GetProfile(userID int) (Profile, error)
	ValidLogin(username, password string) (userID int, error)
	CreateAccount(username, password string) (userID int, error)
	UpdateProfile(userID int, opts UpdateProfileOptions) error
	ChangePassword(userID int, oldPass, newPass string) error
}