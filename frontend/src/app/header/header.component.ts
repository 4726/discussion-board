import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GatewayService } from '../gateway.service';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  providers: [GatewayService]
})
export class HeaderComponent implements OnInit {
  searchForm;
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
    this.testInit()
  }

  prodInit() {
    const userID = this.gatewayService.getUserID()
    if (userID != 0) {
      this.signedIn = true
      this.userID = userID
    } else {
      this.signedIn = false
    }
  }

  testInit() {
    const userID: number = 1
    if (userID != 0) {
      this.signedIn = true
      this.userID = userID
    } else {
      this.signedIn = false
    }
  }

  onSearch(formData) {
    this.router.navigate(['search'], {queryParams: {term: formData.term, page: '1'}})
  }

  onProfileClick() {
    this.router.navigate([`profile/${this.userID}`])
  }
}
