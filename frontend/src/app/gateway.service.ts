import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})

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

const httpOptions = {
  headers: new HttpHeaders({
    'Authorization': 'Bearer <token>'
  })
};

export class GatewayService {
  gatewayAddr = '127.0.0.1:14000'

  constructor(private http: HttpClient) { }

  getPosts() {
    return this.http.get<[Post]>(this.gatewayAddr + "/posts")
      .pipe(
        retry(3),
        catchError(this.handleError)
      )
  }

  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      console.error('client side error:', error.error.message);
    } else {
      console.error(`gateway error: ${error.status}`);
    }

    return throwError('Could not get posts');
  }
}
