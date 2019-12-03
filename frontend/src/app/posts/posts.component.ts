import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { MatTableDataSource } from '@angular/material/table';
import {Post, PostComment} from '../interfaces.service';

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
  hasPrevPage: boolean = true;
  dataSource: MatTableDataSource<Post>;

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private route: ActivatedRoute,
    ) {
      this.dataSource = new MatTableDataSource(this.posts);
  }

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      const pageParam = this.route.snapshot.paramMap.get('page')
      this.page = +pageParam
      this.prodInit()
    })    
  }

  prodInit() {
    this.posts = []
    this.gatewayService.getPosts(this.page)
      .subscribe(
        (data: Post[]) => {
          this.posts = data
          this.dataSource.data = this.posts
        },
        error => this.error = error
      );
  }

  nextPage() {
    this.router.navigate([`posts/${this.page + 1}`])
  }

  prevPage() {
    if (this.page == 1) {
      this.router.navigate([`posts/1`])
    } else {
      this.router.navigate([`posts/${this.page - 1}`])
    }
  }

  onRowClick(postID: number) {
    this.router.navigate([`post/${postID}`])
  }
}