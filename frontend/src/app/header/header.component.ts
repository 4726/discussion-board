import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GatewayService } from '../gateway.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  providers: [GatewayService]
})
export class HeaderComponent implements OnInit {
  signedIn: boolean;
  userID: number;

  constructor(
    private router: Router,
    private gatewayService: GatewayService,
  ) { }

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

}
