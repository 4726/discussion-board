import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, throwError, of } from 'rxjs';
import { catchError, retry, map } from 'rxjs/operators';
import { Post, Profile, InterfacesService } from './interfaces.service';

@Injectable()
export class GatewayService {
  gatewayAddr = 'http://100.115.92.200:14000'

  constructor(
    private http: HttpClient, 
    private interfaces: InterfacesService,
  ) {}

  getPosts(page: number, userID: number = 0): Observable<Post[]> {
    return this.http.get<Post[]>(this.gatewayAddr + `/posts?page=${page}&userid=${userID}`)
      .pipe(
        map(resp => this.interfaces.postsFromJSON(resp["posts"])),
        retry(3),
        catchError(this.handleError)
      )
  }

  getPost(postID: number): Observable<Post> {
    return this.http.get<Post>(this.gatewayAddr + `/post/${postID}`)
      .pipe(
        map(resp => this.interfaces.postFromJSON(resp)),
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
        map(resp => resp["post_id"]),
        retry(3),
        catchError(this.handleError)
      )
  }

  login(username: string, password: string): Observable<string> {
    const postData = {
      username: username,
      password: password,
    }
    return this.http.post<string>(this.gatewayAddr + '/login', postData, {})
      .pipe(
        map(resp => resp['jwt']),
        retry(3),
        catchError(this.handleError)
      )
  }

  register(username: string, password: string): Observable<string> {
    const postData = {
      username: username,
      password: password,
    }
    return this.http.post<string>(this.gatewayAddr + '/register', postData, {})
      .pipe(
        map(resp => resp['jwt']),
        retry(3),
        catchError(this.handleError)
      )
  }

  deletePost(postID: number) {
    const postData = {
      post_id: postID
    }

    this.http.post(this.gatewayAddr + `/post/delete`, postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  likePost(postID: number): Observable<number> {
    const postData = {
      id: postID,
    }

    return this.http.post<number>(this.gatewayAddr + '/post/like', postData, {})
      .pipe(
        map(resp => resp["total"]),
        retry(3),
        catchError(this.handleError)
      )
  }

  unlikePost(postID: number): Observable<number> {
    const postData = {
      id: postID,
    }

    return this.http.post<number>(this.gatewayAddr + '/post/unlike', postData, {})
      .pipe(
        map(resp => resp["total"]),
        retry(3),
        catchError(this.handleError)
      )
  }

  addComment(postID: number, body: string, parentID: number = 0) {
    const postData = {
      post_id: postID,
      body: body,
      parent_id: parentID,
    }

    this.http.post(this.gatewayAddr + '/comment/create', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  likeComment(commentID: number): Observable<number> {
    const postData = {
      id: commentID,
    }

    return this.http.post<number>(this.gatewayAddr + '/comment/like', postData, {})
      .pipe(
        map(resp => resp["total"]),
        retry(3),
        catchError(this.handleError)
      )
  }

  unlikeComment(commentID: number): Observable<number> {
    const postData = {
      id: commentID,
    }

    return this.http.post<number>(this.gatewayAddr + '/comment/unlike', postData, {})
      .pipe(
        map(resp => resp["total"]),
        retry(3),
        catchError(this.handleError)
      )
  }

  clearComment(commentID: number) {
    const postData = {
      comment_id: commentID,
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
        map(resp => this.interfaces.postsFromJSON(resp["posts"])),
        retry(3),
        catchError(this.handleError)
      )
  }

  changePassword(oldPass: string, newPass: string) {
    const postData = {
      old_pass: oldPass,
      new_pass: newPass,
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
        map(resp => this.interfaces.profileFromJSON(resp)),
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
            return resp['user_id']
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
