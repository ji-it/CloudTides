import { Inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DOCUMENT } from '@angular/common';
import { isEmpty } from 'lodash';
import { BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { base } from '@tide-environments/base'

@Injectable()
export class LoginService {

  constructor(
    private readonly http: HttpClient,
    @Inject(DOCUMENT) private readonly document: Document,
  ) { }

  readonly session$ = new BehaviorSubject<UserInfo>({} as any);

  login(
    username = '',
    password = '',
  ) {
    return this.http.post<UserInfo>(`${base.apiPrefix}/session`, { username, password }).pipe(
      tap(userInfo => {
        this.session$.next(userInfo);
      }),
    );
  }

  current() {
    return this.http.get<UserInfo>(`${base.apiPrefix}/session`).pipe(
      tap(userInfo => {
        this.session$.next(userInfo);
      }),
    );
  }

  logout() {
    return this.http.post(`${base.apiPrefix}/session`, {}).pipe(
      tap(() => {
        this.document.location.href = '/login';
      }),
    );
  }

  get session() {
    return this.session$.value;
  }

  get hasLoggedIn() {
    return isEmpty(this.session) === false;
  }
}

export interface UserInfo {
  id: string;
  name: string;
  username: string;
  email: string;
  phone: string;
  website: string;
  admin: boolean;
}
