import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { ActivatedRoute, ParamMap } from "@angular/router";
import {Post, PostComment} from '../interfaces.service';

@Component({
  selector: 'app-get-post',
  templateUrl: './get-post.component.html',
  styleUrls: ['./get-post.component.scss'],
  providers: [GatewayService]
})
export class GetPostComponent implements OnInit {
  post: Post
  error: string;

  constructor(
    private gatewayService: GatewayService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      const param = this.route.snapshot.paramMap.get('postID')
      this.setPost(+param);
    })
  }

  setPost(postID: number) {
    this.gatewayService.getPost(postID)
      .subscribe(
        (data: Post) => this.post = data,
        error => this.error = error
      );
  }

  onLikePost() {
    this.gatewayService.likePost(this.post.ID)
      .subscribe(
        (data: number) => {
          this.post.Likes = data
          this.post.HasLike = true
        },
        error => this.error = error
      )
  }

  onUnlikePost() {
    this.gatewayService.unlikePost(this.post.ID)
      .subscribe(
        (data: number) => {
          this.post.Likes = data
          this.post.HasLike = false
        },
        error => this.error = error
      )
  }

  onLikeComment(commentID: number) {
    this.gatewayService.likeComment(commentID)
    .subscribe(
      (data: number) => {
        let comment = this.findComment(commentID)
        comment.Likes = data
        comment.HasLike = true
      },
      error => this.error = error
    )
  }

  onUnlikeComment(commentID: number) {
    this.gatewayService.unlikeComment(commentID)
      .subscribe(
        (data: number) => {
          let comment = this.findComment(commentID)
          comment.Likes = data
          comment.HasLike = false
        },
        error => this.error = error
      )
  }

  private findComment(commentID: number): PostComment {
    return this.post.Comments.find(v => v.ID == commentID)
  }
}
