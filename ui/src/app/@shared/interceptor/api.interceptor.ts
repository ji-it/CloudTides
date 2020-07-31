import { Provider } from '@angular/core';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs';

import { LoginService } from '../../login/login.service';

function APIInspectorFactory(
  loginService: LoginService,
  ) {
  return new APIInterceptor(loginService);
}

@Injectable()
class APIInterceptor implements HttpInterceptor {

  constructor(
    private readonly loginService: LoginService,
  ) {
  }

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler,
  ): Observable<HttpEvent<any>> {

    const session = this.loginService.session;

    if(session && session.token) {
      const params = req.params;

      let headers = req.headers;
      // Add csp-auth-token in header except when url starts with ${VMWARE_KB_URL_ON_CONTEXTUAL}
      // should not append Token for /orginfo call.
      if (!req.url.endsWith('user/login')) {
        headers = req.headers.set('Authorization', `Bearer ${session.token || ''}`);
      }

      req = req.clone({
        headers,
        params,
      });
    }

    return next.handle(req);
  }
}

export const interceptorProviders: Provider[] = [
  {
    provide: HTTP_INTERCEPTORS,
    useFactory: APIInspectorFactory,
    multi: true,
    deps: [
      LoginService,
    ],
  },
];
