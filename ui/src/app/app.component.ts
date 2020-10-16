import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { PRODUCT_NAME } from '@tide-config/const';

import { LoginService, UserInfo } from './login/login.service';

@Component({
  selector: 'tide-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {

  constructor(
    private readonly loginService: LoginService,
    private readonly router: Router,
  ) {

  }

  readonly vo = {
    title: PRODUCT_NAME,
    user$: this.loginService.session$,
  };

  signOut() {
    this.loginService.logout().subscribe();
  }

  inLoginPage() {
    return this.router.url.indexOf('/login') === 0;
  }

  ngOnInit() {

  }
}
