import { Injectable } from '@angular/core';

export interface Post {
  ID: number;
  UserID: number;
  Title: string;
  Body: string;
  Likes: number;
  CreatedAt: string;
  UpdatedAt: string;
  Comments: PostComment[];
  HasLike: boolean;
}

export interface PostComment {
  ID: number;
  PostID: number;
  ParentID: number;
  UserID: number;
  Body: string;
  CreatedAt: string;
  Likes: number;
  HasLike: boolean;
}

export interface Profile {
  UserID: number;
  Username: string;
  Bio: string;
  AvatarID: string;
  IsMine: boolean;
}

@Injectable()
export class InterfacesService {

  constructor() { }

  postFromJSON(json: any): Post {
    var post = {} as Post
    post.ID = json["id"]
    post.UserID = json["user_id"]
    post.Title = json["title"]
    post.Body = json["body"]
    post.Likes = json["likes"]
    post.CreatedAt = json["created_at"]
    post.UpdatedAt = json["updated_at"]
    post.Comments = this.commentsFromJSON(json["comments"])
    post.HasLike = json["has_like"]
    return post
  }

  commentFromJSON(json: any): PostComment {
    var comment = {} as PostComment
    comment.ID = json["id"]
    comment.PostID = json["post_id"]
    comment.ParentID = json["parent_id"]
    comment.UserID = json["user_id"]
    comment.Body = json["body"]
    comment.CreatedAt = json["created_at"]
    comment.Likes = json["likes"]
    comment.HasLike = json["has_like"]
    return comment
  }

  postsFromJSON(json: any[]): Post[] {
    return json.map(v => {
      return this.postFromJSON(v)
    })
  }

  commentsFromJSON(json: any[]): PostComment[] {
    return json.map(v => {
      return this.commentFromJSON(v)
    })
  }

  profileFromJSON(json: any): Profile {
    var profile = {} as Profile
    profile.UserID = json["user_id"]
    profile.Username = json["username"]
    profile.Bio = json["bio"]
    profile.AvatarID = json["avatar_id"]
    profile.IsMine = json["is_mine"]
    return profile
  }
}