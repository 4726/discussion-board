import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

export interface Post {
  ID: number;
  UserID: number;
  Title: string;
  Body: string;
  Likes: number;
  CreatedAt: string;
  UpdatedAt: string;
  Comments: PostComment;
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

@Injectable()
export class GatewayService {
  gatewayAddr = '127.0.0.1:14000'

  constructor(private http: HttpClient) { }

  getPosts() {
    return this.http.get<[Post]>(this.gatewayAddr + '/posts')
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

    const resp = this.http.post(this.gatewayAddr + '/login', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
    return resp['jwt']
  }

  register(username: string, password: string): Observable<string> {
    const postData = {
      username: username,
      password: password,
    }

    const resp = this.http.post(this.gatewayAddr + '/register', postData, {})
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
    return resp['jwt']
  }

  validJWT(): boolean {
    let statusCode
    this.http.get(this.gatewayAddr + '/login', {observe: 'response'})
      .subscribe(resp => {
          statusCode = resp.status
      })
    return statusCode == 400
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
