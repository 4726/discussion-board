import { Component, OnInit } from '@angular/core';
import { GatewayService, Post, PostComment } from '../gateway.service';

@Component({
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss']
})

export class PostsComponent implements OnInit {
  posts: [Post];
  error: string;

  constructor(private gatewayService: GatewayService) { }

  ngOnInit() {
  }

  showPosts() {
    this.gatewayService.getPosts()
      .subscribe(
        (data: [Post]) => this.posts = {...data},
        error => this.error = error
      );
  }

}