import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { PRODUCT_NAME } from '@tide-config/const';

import { LoginService, UserInfo } from './login/login.service';
import { TranslateService } from '@ngx-translate/core';
import { I18nService } from '@tide-shared/service/i18n';
import { Observable, Subject } from 'rxjs';
import { RegisterService } from './register/register.service';

@Component({
  selector: 'tide-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {

  constructor(
    readonly loginService: LoginService,
    readonly registerService: RegisterService,
    private readonly router: Router,
    translate: TranslateService,
    public readonly i18nService: I18nService,
  ) {
  }

  readonly vo = {
    title: PRODUCT_NAME,
  };

  subject = new Subject();

  signOut() {
    this.loginService.logout();
  }

  ngOnInit() {

  }
}
