import { Component, OnInit, OnDestroy } from '@angular/core';
import { EMPTY, Subject } from 'rxjs';

import { LoginService } from './login.service';
import { catchError, switchMap, tap } from 'rxjs/operators';
import { Router } from '@angular/router';

@Component({
  selector: 'cp-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit, OnDestroy {

  constructor(
    public readonly loginService: LoginService,
    private readonly router: Router,
  ) {
  }

  readonly vo = {

    model: {
      username: '',
      password: '',
    },

    submitting: false,
    loginError: '',

  };

  private readonly submit$ = new Subject<Credential>();

  private readonly submit$$ = this.submit$.asObservable()
    .pipe(
      tap(() => {
        this.vo.submitting = true;
        this.vo.loginError = '';
      }),
      switchMap(({ username, password }) => {

        return this.loginService
          .login(username.trim(), password)
          .pipe(
            tap(() => {
              this.vo.submitting = false;
            }),
            catchError((error, source) => {
              this.vo.submitting = false;
              this.vo.loginError = error.message;

              return EMPTY as typeof source;
            }),
          );
      }),
    )
    .subscribe(res => {
      // this.document.location.href = '/'
      this.router.navigate(['/']);
    })
  ;

  onSubmit({ username = '', password = '' }: Credential) {
    this.submit$.next({ username, password });
  }

  ngOnInit() {
    if (this.loginService.hasLoggedIn) {
      // this.document.location.href = '/';
      this.router.navigate(['/']);
    }
  }

  ngOnDestroy() {
    this.submit$$.unsubscribe();
  }
}

interface Credential {
  username: string;
  password: string;
}
