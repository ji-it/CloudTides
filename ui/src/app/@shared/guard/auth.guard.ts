import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import {
  CanLoad,
  CanActivate,
  CanActivateChild,
  Route,
  RouterStateSnapshot,
  ActivatedRouteSnapshot,
} from '@angular/router';
import { get } from 'lodash';
import { tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';

import { LoginService } from '../../login/login.service';
import { RouterData } from '../../app-routing.module';
import { LOGIN_PATH } from '@tide-config/path';

@Injectable()
export class AuthGuard implements CanActivate, CanActivateChild, CanLoad {
  constructor(
    private readonly loginService: LoginService,
    private readonly router: Router,
  ) {
  }

  canLoad(route: Route) {
    return this.check();
  }

  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    const data = get(next, 'data', {}) as RouterData;

    if (data.anonymous === true) {
      return true;
    }

    return this.check().pipe(
      tap(async valid => {
        if (valid) {
          return;
        } else {
          await this.router.navigate([LOGIN_PATH]);
          return;
        }
      }),
    );
  }

  canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    return this.canActivate(childRoute, state);
  }

  private check(): Observable<boolean> {
    if (this.loginService.hasLoggedIn) {
      return of(true);
    } else {
      return of(false);
    }
  }
}
