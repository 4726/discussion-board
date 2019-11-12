
	LikePost(postID, userID int) error  //todo: userID because might already liked
	UnlikePost(postID, userID int) error //todo
	LikeComment(commentID, userID int) //todo
	UnlikeComment(commentID, userID int) //todo
	GetMultiple(postIDs []int) ([]Post, error) //todo