import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { GatewayService } from './gateway.service';
import { catchError, map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})

//user should not be able to view /login route if already logged in
export class NoAuthGuard implements CanActivate {

  constructor(
    private gatewayService: GatewayService,
    private router: Router,
  ){}

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
      return this.gatewayService.validJWT()
        .pipe(
          map(valid => {
            if (valid) {
              this.router.navigate(['home'])
              return false
            } else {
              return true
            }
          })
        )
  }
  
}
