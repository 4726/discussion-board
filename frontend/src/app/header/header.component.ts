import { Component, OnInit } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { GatewayService } from '../gateway.service';
import { FormBuilder, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  providers: [GatewayService]
})
export class HeaderComponent implements OnInit {
  searchForm: FormGroup;
  signedIn: boolean;
  userID: number;

  constructor(
    private router: Router,
    private gatewayService: GatewayService,
    private formBuilder: FormBuilder,
  ) { 
    this.searchForm = this.formBuilder.group({
      term: ''
    });
  }

  ngOnInit() {
    this.prodInit()
    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        this.prodInit()
      }
    })
  }

  prodInit() {
    this.gatewayService.getUserID()
      .subscribe(
        res => {
          this.signedIn = true
          this.userID = res
        },
        err => {
          this.signedIn = false
        }
    )
  }

  onSearch(formData) {
    this.router.navigate(['search'], {queryParams: {term: formData.term, page: '1'}})
  }

  onProfileClick() {
    this.router.navigate([`profile/${this.userID}`])
  }

  onLoginClick() {
    this.router.navigate(['login'])
  }

  onRegisterClick() {
    this.router.navigate(['register'])
  }

  onLogoutClick() {
    localStorage.removeItem('jwt');
    this.signedIn = false
    this.router.navigate(['home'])
  }
}
