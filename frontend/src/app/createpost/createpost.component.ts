import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GatewayService } from '../gateway.service';
import { FormBuilder, ReactiveFormsModule, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-post',
  templateUrl: './createpost.component.html',
  styleUrls: ['./createpost.component.scss'],
  providers: [GatewayService]
})
export class CreatePostComponent implements OnInit {
  createForm: FormGroup;
  error: string

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private formBuilder: FormBuilder,
  ) {
    this.createForm = this.formBuilder.group({
      title: '',
      body: ''
    });
  }

  ngOnInit() {
  }

  onSubmit(postData) {
    this.gatewayService.createPost(postData.title, postData.body)
      .subscribe(
        (data: number) => this.router.navigate([`/post/${data}`]),
        error => this.error = error
      );
  }

}
