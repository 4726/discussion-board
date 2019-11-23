import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, throwError, of } from 'rxjs';
import { catchError, retry, map } from 'rxjs/operators';

export interface Post {
  ID: number;
  UserID: number;
  Title: string;
  Body: string;
  Likes: number;
  CreatedAt: string;
  UpdatedAt: string;
  Comments: PostComment[];
}

export interface PostComment {
  ID: number;
  PostID: number;
  ParentID: number;
  UserID: number;
  Body: string;
  CreatedAt: string;
  Likes: number;
}

export interface Profile {
  UserID: number;
  Username: string;
  Bio: string;
  AvatarID: string;
  IsMine: boolean;
}

@Injectable()
export class GatewayService {
  gatewayAddr = 'http://100.115.92.200:14000'

  constructor(private http: HttpClient) { }

  getPosts(page: number, userID: number = 0): Observable<Post[]> {
    return this.http.get<Post[]>(this.gatewayAddr + `/posts?page=${page}&userid=${userID}`)
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  getPost(postID: number): Observable<Post> {
    return this.http.get<Post>(this.gatewayAddr + `/post/${postID}`)
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  createPost(title: string, body: string): Observable<number> {
    const post = {
      title: title,
      body: body,
    }
    return this.http.post<number>(this.gatewayAddr + '/post', post, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  login(username: string, password: string): Observable<string> {
    const postData = {
      username: username,
      password: password,
    }

    return this.http.post(this.gatewayAddr + '/login', postData, {})
      .pipe(
        map(
          res => res['jwt']
        ),
        retry(3),
        catchError(this.handleError)
      )
  }

  register(username: string, password: string): Observable<string> {
    const postData = {
      username: username,
      password: password,
    }


    return this.http.post(this.gatewayAddr + '/register', postData, {})
      .pipe(
        map(
          res => res['jwt']
        ),
        retry(3),
        catchError(this.handleError)
      )
  }

  deletePost(postID: number) {
    this.http.delete(this.gatewayAddr + `/post/${postID}`, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  likePost(postID: number) {
    const postData = {
      postID: postID,
    }

    this.http.post(this.gatewayAddr + '/post/like', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  unlikePost(postID: number) {
    const postData = {
      postID: postID,
    }

    this.http.post(this.gatewayAddr + '/post/unlike', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  addComment(postID: number, body: string, parentID: number = 0) {
    const postData = {
      postID: postID,
      body: body,
      parentID: parentID,
    }

    this.http.post(this.gatewayAddr + '/comment/create', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  likeComment(commentID: number) {
    const postData = {
      commentID: commentID,
    }

    this.http.post(this.gatewayAddr + '/comment/like', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  unlikeComment(commentID: number) {
    const postData = {
      commentID: commentID,
    }

    this.http.post(this.gatewayAddr + '/comment/unlike', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  clearComment(commentID: number) {
    const postData = {
      commentID: commentID,
    }

    this.http.post(this.gatewayAddr + '/comment/clear', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  search(term: string, page: number): Observable<Post[]> {
    return this.http.get<Post[]>(this.gatewayAddr + `/search?term=${term}&page=${page}`, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  changePassword(oldPass: string, newPass: string) {
    const postData = {
      oldPass: oldPass,
      newPass: newPass,
    }

    this.http.post(this.gatewayAddr + '/changepassword', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  getProfile(userID: number): Observable<Profile> {
    return this.http.get<Profile>(this.gatewayAddr + `/profile/${userID}`, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  updateProfile(bio: string) {
    const postData = {
      bio: bio,
    }

    this.http.post(this.gatewayAddr + '/profile/update', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  validJWT(): Observable<boolean> {
    return this.getUserID()
      .pipe(
        map(
          resp => {
            return true
          }
        ),
        catchError((error) => {
          return of(false)
        })
      )
  }

  getUserID(): Observable<number> {
    return this.http.get(this.gatewayAddr + '/userid')
      .pipe(
        map(
          resp => {
            console.log(resp)
            return resp['UserID']
          }
        )
      )
  }

  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      console.error(`client side error: ${error.error.message}`);
    } else {
      console.error(`gateway error: ${error.status}`);
    }

    return throwError('Could not get posts');
  }
}
