import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { TemplateService, Item, ItemVM } from '../template.service';
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
  VMlist$: Observable<ItemVM[]> = of([]);
  TemplateList: Object = {};
  //ResList: Object ={};
  opened = false;
  VMopened = false;
  displayOpened = false;
  TemplateID = 1;
  refreshInterval: number;

  async save() {
    await this.refreshList();
  }

  open() {
    this.opened = true;
  }

  cancel() {
    this.opened = false;
    this.VMopened = false;
    this.displayOpened = false;
  }

  async ngOnInit() {
    await this.refreshList();
  }

  async refreshList() {
    this.list$ = of(await this.templateService.getList());
    this.refreshInterval = window.setInterval(async () => {
      this.list$ = of(await this.templateService.getList());
    }, VENDOR_USAGE_REFRESH_PERIOD);
    //this.ResList = Object(await this.templateService.getResList());
    this.TemplateList = Object(await this.templateService.getTemplateList());
  }

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

  async add(id: number) {
    this.TemplateID = id;
    this.VMopened = true;
  }

  async publish(id: number) {
  }

  async displayVM(id: number) {
    this.TemplateID = id;
    this.displayOpened = true;
    this.VMlist$ = of(await this.templateService.getVMList(id))
  }

  async delete(id: string) {
    await this.templateService.removeItem(id).then(() => {
      this.vo.alertText = `Successfully delete vendor with id ${id}`;
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
    this.refreshList();
  }
}
