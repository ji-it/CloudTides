import { Inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { REGISTER_API_URL, REGISTER_PATH } from '@tide-config/path';
import { environment } from '@tide-environments/environment';
import { tap } from 'rxjs/operators';
import { DOCUMENT } from '@angular/common';

@Injectable()
export class RegisterService {

  constructor(
    private readonly http: HttpClient,
    @Inject(DOCUMENT) private readonly document: Document,
  ) {

  }

  register(
    username = '',
    password = '',
    companyName = '',
    phone = '',
    email = '',
    priority = Priority.LOW,
  ) {
    return this.http.post<RegisterResult>(environment.apiPrefix + REGISTER_API_URL,
      { username, password, priority, companyName, phone, email }).pipe(
      tap(val => {

      }),
    );
  }

  inRegisterPage() {
    return this.document.location.pathname === REGISTER_PATH;
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

