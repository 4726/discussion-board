import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
  providers: [GatewayService]
})

export class LoginComponent implements OnInit {
  loginForm;
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
        (data: string) => {
          localStorage.setItem('jwt', data);
          this.router.navigate(['/home']);
        },
        error => this.error = error
      );
  }

}
