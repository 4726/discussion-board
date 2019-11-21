import { Component, OnInit } from '@angular/core';
import { GatewayService, Profile } from '../gateway.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

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
      // this.prodInit()
      this.testInit()
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

  testInit() {
    const p = {} as Profile
    p.UserID = this.userID
    p.Username = 'my_username'
    p.Bio = 'hello world'
    p.AvatarID = ''

    if (this.userID == 1) {
      p.IsMine = true
    } else {
      p.IsMine = false
    }
    this.profile = p
    this.isMine = p.IsMine
  }

  onEditProfile() {

  }

  onViewPosts() {
    this.router.navigate([`profile/${this.userID}/posts/1`])
  }
}
