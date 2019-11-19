import { Component, OnInit } from '@angular/core';
import { GatewayService, Post, PostComment } from '../gateway.service';
import { ActivatedRoute } from "@angular/router";

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
    const param = this.route.snapshot.paramMap.get('postID')
    if (!param) {
      this.getPost(+param);
    } else {
      //404 page
    }
  }

  getPost(postID: number) {
    this.gatewayService.getPost(postID)
      .subscribe(
        (data: Post) => this.post = data,
        error => this.error = error
      );
  }

}
