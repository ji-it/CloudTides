import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { REGISTER_PATH } from '@tide-config/path';
import { base } from '@tide-environments/base';
import { tap } from 'rxjs/operators';
import { ServerUserInfo, UserInfo } from '../login/login.service';

@Injectable()
export class RegisterService {

  constructor(
    private readonly http: HttpClient,
  ) {

  }

  register(
    username = '',
    password = '',
    companyName = '',
    priority = Priority.LOW,
  ) {
    return this.http.post<RegisterResult>(base.apiPrefix + REGISTER_PATH,
      { username, password, priority, companyName }).pipe(
      tap(val => {

      }),
    );
  }
}

export interface RegisterResult {
  userinfo: {
    username: string,
    priority: string,
    password: string,
  }
}

enum Priority {
  HIGH = 'High',
  MEDIUM = 'Medium',
  LOW = 'Low',
}

