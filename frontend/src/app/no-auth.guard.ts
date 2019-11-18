import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router } from '@angular/router';
import { Observable } from 'rxjs';
// import { GatewayService } from './gateway.service';

@Injectable({
  providedIn: 'root'
})

//user should not be able to view /login route if already logged in
export class NoAuthGuard implements CanActivate {

  constructor(
    // private gatewayService: GatewayService,
    private router: Router,
  ){}

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    // const loggedIn = this.gatewayService.validJWT()
    const loggedIn = true
    if (loggedIn) {
      this.router.navigate(['/']);
      return false
    } else {
      return true
    }
  }
  
}
