import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import {Profile} from '../interfaces.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
  providers: [GatewayService]
})
export class ProfileComponent implements OnInit {
  userID: number
  error: string
  profile: Profile
  isMine: boolean

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private route: ActivatedRoute,
    ) { }

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      const userIDParam = this.route.snapshot.paramMap.get('userid')
      this.userID = +userIDParam
      this.prodInit()
    })    
  }

  prodInit() {
    this.gatewayService.getProfile(this.userID)
    .subscribe(
      (data: Profile) => {
        this.profile = {...data}
      },
      error => this.error = error
    );
  }

  onEditProfile() {

  }

  onViewPosts() {
    this.router.navigate([`profile/${this.userID}/posts/1`])
  }
}
