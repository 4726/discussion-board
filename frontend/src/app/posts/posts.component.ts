import { Component, OnInit } from '@angular/core';
import { GatewayService, Post, PostComment } from '../gateway.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss'],
  providers: [GatewayService]
})

export class PostsComponent implements OnInit {
  posts: Post[] = [];
  error: string;
  displayedColumns: string[] = ['title', 'userid', 'likes', 'updatedat']//html

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    ) {
  }

  ngOnInit() {
    // this.showPosts()
    const p = {} as Post
    p.ID = 123
    p.UserID = 321
    p.Title = 'hello world'
    p.Likes = 1
    p.UpdatedAt = '1 hour ago'

    this.posts.push(p)
  }

  showPosts() {
    this.gatewayService.getPosts()
      .subscribe(
        (data: [Post]) => this.posts = {...data},
        error => this.error = error
      );
  }

}