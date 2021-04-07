import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { TemplateService, Item } from '../template.service';
import { TranslateService } from '@ngx-translate/core';
import { LoginService } from 'src/app/login/login.service';
import { NOTIFICATION_EXIST_TIME, VENDOR_USAGE_REFRESH_PERIOD } from '@tide-shared/config/const';

@Component({
  selector: 'tide-template-list',
  templateUrl: './template-list.component.html',
  styleUrls: ['./template-list.component.scss'],
})
export class TemplateListComponent implements OnInit, OnDestroy{

  constructor(
    public templateService: TemplateService,
    public readonly translate: TranslateService,
    public readonly loginService: LoginService,
  ) { }

  readonly vo = {
    alertType: '',
    alertText: '',
  };

  async resetAlert(time?: number) {
    window.setTimeout(() => {
      this.vo.alertText = '';
    }, time || NOTIFICATION_EXIST_TIME);
  }

  list$: Observable<Item[]> = of([]);
  opened = false;
  refreshInterval: number;

  async save() {
    await this.refreshList();
  }

  cancel() {
    this.opened = false;
  }

  async ngOnInit() {
    await this.refreshList;
  }

  async refreshList() {
    this.list$ = of(await this.templateService.getList());
    this.refreshInterval = window.setInterval(async () => {
      this.list$ = of(await this.templateService.getList());
    }, VENDOR_USAGE_REFRESH_PERIOD);
  }

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

  async add() {
    this.list$ = of(await this.templateService.getList());
  }
}
