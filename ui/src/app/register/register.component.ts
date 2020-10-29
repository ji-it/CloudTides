import { Component, OnDestroy, OnInit } from '@angular/core';
import { LoginService } from '../login/login.service';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { I18nService } from '@tide-shared/service/i18n';
import { EMPTY, Subject } from 'rxjs';
import { catchError, switchMap, tap } from 'rxjs/operators';
import { RegisterService } from './register.service';
import { FormBuilder, Validators, FormGroup, FormControl, ValidatorFn } from '@angular/forms';

function passwordMatchValidator(password: string): ValidatorFn {
  return (control: FormControl) => {
    if (!control || !control.parent) {
      return null;
    }
    return control.parent.get(password).value === control.value ? null : { mismatch: true };
  };
}

@Component({
  selector: 'tide-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent implements OnInit, OnDestroy {

  registerForm: FormGroup;

  constructor(
    public readonly loginService: LoginService,
    public readonly registerService: RegisterService,
    private readonly router: Router,
    public readonly translate: TranslateService,
    public readonly i18nService: I18nService,
    private fb: FormBuilder,
  ) {
    this.registerForm = this.fb.group({
      username: [
        '', [
          Validators.required,
          Validators.minLength(4),
          Validators.maxLength(12),
        ],
      ],
      password: [
        '', [
          Validators.required,
          Validators.minLength(6),
          Validators.maxLength(16),
        ]],
      confirmPassword: [
        '', [
          Validators.required,
          Validators.minLength(6),
          Validators.maxLength(16),
          passwordMatchValidator('password'),
        ]],
      companyName: [
        '', [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(10),
        ],
      ],
    });
  }

  readonly vo = {
    submitting: false,
    loginError: '',
  };

  private readonly submit$ = new Subject<RegisterForm>();

  private readonly submit$$ = this.submit$.asObservable()
    .pipe(
      tap(() => {
        this.vo.submitting = true;
        this.vo.loginError = '';
      }),
      switchMap(({ username, password, companyName }) => {
        return this.registerService
          .register(username, password, companyName)
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
      this.router.navigate(['/login']);
    })
  ;

  onSubmit() {
    this.submit$.next(this.registerForm.value);
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

interface RegisterForm {
  username: string;
  password: string;
  companyName: string,
}
