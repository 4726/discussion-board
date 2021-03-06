import { Injectable } from '@angular/core';
import {
  HttpEvent, HttpInterceptor, HttpHandler, HttpRequest, HttpHeaders
} from '@angular/common/http';

import { Observable } from 'rxjs';

@Injectable()
export class JWTInterceptor implements HttpInterceptor {

    constructor(){}

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        const token = localStorage.getItem("jwt")
        console.log('token: ' + token)
        if (token != null) {
            const updatedReq = req.clone({
                setHeaders: {Authorization: `Bearer ${token}`}
            });
            return next.handle(updatedReq)
        }
        return next.handle(req);
    }
}

