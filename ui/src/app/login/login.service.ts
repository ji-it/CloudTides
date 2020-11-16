import { Inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DOCUMENT } from '@angular/common';
import { BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { base } from '@tide-environments/base';
import { LOGIN_PATH, PROFILE_PATH } from '@tide-config/path';
import { LOCAL_STORAGE_KEY } from '@tide-config/const';
import { Router } from '@angular/router';

@Injectable()
export class LoginService {

  constructor(
    private readonly http: HttpClient,
    private readonly router: Router,
    @Inject(DOCUMENT) private readonly document: Document,
  ) {
    this.loginNavigate();
  }

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

  loginNavigate() {
    if (this.hasLoggedIn) {
      this.current().subscribe(
        () => {},
        error => {
          switch (error.status) {
            case 401:
              this.logout();
              break;
            default:
              break;
          }
        });
    } else {
      this.router.navigate(['/']);
    }
  }

  current() {
    return this.http.get<UserInfo>(base.apiPrefix + PROFILE_PATH, {
      headers: {
        Authorization: `Bearer ${this.token}`,
      },
    }).pipe(
      tap(userInfo => {
        // todo: handle the condition when userInfo is empty (backend part)
        this.session$.next(userInfo);
      }),
    );
  }

  logout() {
    this.removeToken();
    this.router.navigate(['/']);
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

