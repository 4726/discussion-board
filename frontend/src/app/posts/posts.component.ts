import { Component, OnInit } from '@angular/core';
import { GatewayService, Post, PostComment } from '../gateway.service';
import { Router, ActivatedRoute } from '@angular/router';
import { MatTableDataSource } from '@angular/material/table';

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
  page: number = 1;
  hasPrevPage: boolean = false;
  dataSource: MatTableDataSource<Post>;

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private route: ActivatedRoute,
    ) {
      this.dataSource = new MatTableDataSource(this.posts);
  }

  ngOnInit() {
    // this.showPosts()
    this.testInit()
  }

  showPosts() {
    this.gatewayService.getPosts(this.page)
      .subscribe(
        (data: [Post]) => this.posts = {...data},
        error => this.error = error
      );
  }

  testInit() {
    this.posts = [];
    if (this.page == 1) {
      const p = {} as Post
      p.ID = 123
      p.UserID = 321
      p.Title = 'hello world'
      p.Likes = 1
      p.UpdatedAt = '1 hour ago'
  
      this.posts.push(p)
  
      const p2 = {} as Post
      p2.ID = 124
      p2.UserID = 140
      p2.Title = 'my title'
      p2.Likes = 20
      p2.UpdatedAt = '5 hours ago'
  
      this.posts.push(p2)

      this.hasPrevPage = false
    } else if (this.page == 2) {
      const p = {} as Post
      p.ID = 222
      p.UserID = 12
      p.Title = 'testing'
      p.Likes = 13
      p.UpdatedAt = '1 day ago'
  
      this.posts.push(p)
  
      const p2 = {} as Post
      p2.ID = 333
      p2.UserID = 1
      p2.Title = 'good'
      p2.Likes = 200
      p2.UpdatedAt = '5 days ago'
  
      this.posts.push(p2)

      this.hasPrevPage = true
    }
    this.dataSource.data = this.posts //refresh table
  }

  nextPage() {
    this.page = this.page + 1
    this.testInit()
  }

  prevPage() {
    this.page = this.page - 1
    this.testInit()
  }


}