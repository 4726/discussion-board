import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { MatTableDataSource } from '@angular/material/table';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import {Post} from '../interfaces.service';

@Component({
  selector: 'app-user-posts',
  templateUrl: './user-posts.component.html',
  styleUrls: ['./user-posts.component.scss'],
  providers: [GatewayService]
})
export class UserPostsComponent implements OnInit {
  posts: Post[] = [];
  error: string;
  displayedColumns: string[] = ['title', 'userid', 'likes', 'updatedat']//html
  page: number = 1;
  hasPrevPage: boolean = true;
  dataSource: MatTableDataSource<Post>;
  userID: number;

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
      const userIDParam = this.route.snapshot.paramMap.get('userid')
      this.userID = +userIDParam
      this.prodInit()
    })    
  }

  prodInit() {
    this.posts = []
    this.gatewayService.getPosts(this.page, this.userID)
      .subscribe(
        (data: Post[]) => {
          this.posts = data
          this.dataSource.data = this.posts
        },
        error => this.error = error
      );
  }

  nextPage() {
    this.router.navigate([`profile/${this.userID}/posts/${this.page + 1}`])
  }

  prevPage() {
    if (this.page == 1) {
      this.router.navigate([`profile/${this.userID}/posts/1`])
    } else {
      this.router.navigate([`profile/${this.userID}/posts/${this.page - 1}`])
    }
  }

  onRowClick(postID: number) {
    this.router.navigate([`post/${postID}`])
  }

}
