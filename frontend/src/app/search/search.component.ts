import { Component, OnInit } from '@angular/core';
import { GatewayService } from '../gateway.service';
import { ActivatedRoute, Router, ParamMap } from '@angular/router';
import { MatTableDataSource } from '@angular/material/table';
import {Post} from '../interfaces.service';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
  providers: [GatewayService]
})
export class SearchComponent implements OnInit {
  posts: Post[] = [];
  error: string;
  displayedColumns: string[] = ['title', 'userid', 'likes', 'updatedat']//html
  page: number = 1;
  term: string = '';
  hasPrevPage: boolean = true;
  dataSource: MatTableDataSource<Post>;

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
    private route: ActivatedRoute,
    ) {
      this.dataSource = new MatTableDataSource(this.posts);
  }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      const pageParam = params['page']
      this.page = +pageParam
      this.term = params['term']
      this.prodInit()
    })  
  }

  prodInit() {
    this.gatewayService.search(this.term, this.page)
  }

  nextPage() {
    this.router.navigate(['search'], {queryParams: {term: this.term, page: this.page + 1}})
  }

  prevPage() {
    if (this.page == 1) {
      this.router.navigate(['search'], {queryParams: {term: this.term, page: 1}})
    } else {
      this.router.navigate(['search'], {queryParams: {term: this.term, page: this.page - 1}})
    }
  }

}
