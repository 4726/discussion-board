import { HTTP_INTERCEPTORS } from '@angular/common/http';

import { JWTInterceptor } from './jwt';

export const httpInterceptorProviders = [
    {provide: HTTP_INTERCEPTORS, useClass: JWTInterceptor, multi: true},
];
