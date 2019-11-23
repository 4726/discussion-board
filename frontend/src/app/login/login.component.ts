import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
  providers: [GatewayService]
})

export class LoginComponent implements OnInit {
  loginForm: FormGroup;
  error: string

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private formBuilder: FormBuilder,
    ){
      this.loginForm = this.formBuilder.group({
        username: '',
        password: ''
      });
    }

  ngOnInit() {
  }

  onSubmit(postData) {
    this.gatewayService.login(postData.username, postData.password)
    .subscribe(
      res => {
        localStorage.setItem('jwt', res);
        this.router.navigate(['/home']);
      },
      err => {
        this.router.navigate(['/home']);
      }
    )
  }

}
