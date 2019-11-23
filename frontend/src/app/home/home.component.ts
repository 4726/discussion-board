import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  providers: [GatewayService]
})
export class HomeComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
