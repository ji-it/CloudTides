import { Inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DOCUMENT } from '@angular/common';
import { isEmpty } from 'lodash';
import { BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { base } from '@tide-environments/base';
import { LOGIN_PATH } from '@tide-config/path';
import { LOCAL_STORAGE_KEY } from '@tide-config/const';

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
    return this.http.post<ServerUserInfo>(base.apiPrefix + LOGIN_PATH, { username, password }).pipe(
      tap(serverUserInfo => {
        this.storeToken(serverUserInfo.token);
        this.session$.next({ ...serverUserInfo.userInfo });
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
    this.removeToken();
    this.document.location.href = '/login';
  }

  storeToken(token: string) {
    localStorage.setItem(LOCAL_STORAGE_KEY.TOKEN, token);
  }

  removeToken() {
    localStorage.removeItem(LOCAL_STORAGE_KEY.TOKEN);
  }

  get session() {
    return this.session$.value;
  }

  get hasLoggedIn() {
    return localStorage.getItem(LOCAL_STORAGE_KEY.TOKEN);
  }

  get token() {
    return localStorage.getItem(LOCAL_STORAGE_KEY.TOKEN);
  }
}

export interface UserInfo {
  username: string;
  password: string;
  priority: string;
  firstName: string;
  lastName: string;
  country: string;
  city: string;
  companyName: string;
  position: boolean;
  email: string,
  phone: string,
}

export interface ServerUserInfo {
  token: string;
  userInfo: UserInfo;
}

