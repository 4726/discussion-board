import { Component, OnInit } from '@angular/core';
import { Post, GatewayService } from '../gateway.service';
import { MatTableDataSource } from '@angular/material/table';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

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
  userID: number;
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
      const userIDParam = this.route.snapshot.paramMap.get('userid')
      this.userID = +userIDParam
      // this.prodInit()
      this.testInit()
    })
  }

  prodInit() {
    this.gatewayService.getPosts(this.page, this.userID)
  }

  testInit() {
    this.posts = [];
    if (this.page == 1) {
      const p = {} as Post
      p.ID = 123
      p.UserID = this.userID
      p.Title = 'hello world'
      p.Likes = 1
      p.UpdatedAt = '1 hour ago'

      this.posts.push(p)

      const p2 = {} as Post
      p2.ID = 124
      p2.UserID = this.userID
      p2.Title = 'my title'
      p2.Likes = 20
      p2.UpdatedAt = '5 hours ago'

      this.posts.push(p2)

      this.hasPrevPage = false
    } else if (this.page == 2) {
      const p = {} as Post
      p.ID = 222
      p.UserID = this.userID
      p.Title = 'testing'
      p.Likes = 13
      p.UpdatedAt = '1 day ago'

      this.posts.push(p)

      const p2 = {} as Post
      p2.ID = 333
      p2.UserID = this.userID
      p2.Title = 'good'
      p2.Likes = 200
      p2.UpdatedAt = '5 days ago'

      this.posts.push(p2)

      this.hasPrevPage = true
    }
    this.dataSource.data = this.posts //refresh table
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

}
