import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
  providers: [GatewayService]
})
export class RegisterComponent implements OnInit {
  registerForm;
  error: string;

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private formBuilder: FormBuilder,
    ){
      this.registerForm = this.formBuilder.group({
        username: '',
        password: ''
      });
    }

  ngOnInit() {
  }

  onSubmit(postData) {
    if (postData.password != postData.password2) {
      this.error = 'password and confirmation password do not match'
      return
    }
    this.gatewayService.register(postData.username, postData.password)
      .subscribe(
        (data: string) => {
          localStorage.setItem('jwt', data);
          this.router.navigate(['/home']);
        },
        error => this.error = error
      );
  }

}
