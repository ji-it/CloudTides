import { Inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DOCUMENT } from '@angular/common';
import { BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { environment } from '@tide-environments/environment';
import { LOGIN_API_URL, LOGIN_PATH, PROFILE_API_URL } from '@tide-config/path';
import { LOCAL_STORAGE_KEY } from '@tide-config/const';
import { Router } from '@angular/router';
import { RegisterService } from '../register/register.service';

@Injectable()
export class LoginService {

  constructor(
    private readonly http: HttpClient,
    private readonly router: Router,
    private readonly registerService: RegisterService,
    @Inject(DOCUMENT) private readonly document: Document,
  ) {
    this.loginNavigate();
  }

  readonly session$ = new BehaviorSubject<UserInfo>({} as any);

  login(
    username = '',
    password = '',
  ) {
    return this.http.post<ServerUserInfo>(environment.apiPrefix + LOGIN_API_URL, { username, password }).pipe(
      tap(serverUserInfo => {
        this.storeToken(serverUserInfo.token);
        this.session$.next({ ...serverUserInfo.userInfo });
      }),
    );
  }

  async loginNavigate() {
    if (this.hasLoggedIn) {
      this.current().subscribe(
        () => {},
        async error => {
          await this.logout();
        });
    } else {
      if (!this.inLoginPage() && !this.registerService.inRegisterPage()) {
        await this.logout();
      }
    }
  }

  current() {
    return this.http.get<any>(environment.apiPrefix + PROFILE_API_URL, {
      headers: {
        Authorization: `Bearer ${this.token}`,
      },
    }).pipe(
      tap(returnMessage => {
        const userInfo = returnMessage.results;
        this.session$.next(userInfo);
      }),
    );
  }

  async logout() {
    this.removeToken();
    await this.router.navigate([LOGIN_PATH]);
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
    return localStorage.getItem(LOCAL_STORAGE_KEY.TOKEN) !== null;
  }

  get token() {
    return localStorage.getItem(LOCAL_STORAGE_KEY.TOKEN);
  }

  inLoginPage() {
    return this.document.location.pathname === LOGIN_PATH;
  }

  inAdminView() {
    return this.session.priority === 'High';
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

