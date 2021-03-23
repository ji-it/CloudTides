import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { Item, ItemPayload, ResourceService } from '../resource.service';
import { TranslateService } from '@ngx-translate/core';
import { NOTIFICATION_EXIST_TIME, RESOURCE_USAGE_REFRESH_PERIOD } from '@tide-config/const';
import { LoginService } from 'src/app/login/login.service';

@Component({
  selector: 'tide-resource-list',
  templateUrl: './resource-list.component.html',
  styleUrls: ['./resource-list.component.scss'],
})
export class ResourceListComponent implements OnInit, OnDestroy {

  constructor(
    private resourceService: ResourceService,
    public readonly translate: TranslateService,
    public readonly loginService: LoginService,
  ) {

  }

  readonly vo = {
    alertType: '',
    alertText: '',
  };

  async contribute(id: string) {
    await this.resourceService.contributeResource(id).then((resp) => {
      if (resp.contributed) {
        this.vo.alertText = `Successfully start contributing Resource with id${id}`;
      } else {
        this.vo.alertText = `Successfully stop contributing Resource with id${id}`;
      }
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
    this.refreshList();
  }

  async activate(id: string) {
    await this.resourceService.activateResource(id).then((resp) => {
      if (resp.activated) {
        this.vo.alertText = `Successfully activate Resource with id ${id}`;
      } else {
        this.vo.alertText = `Successfully deactivate Resource with id ${id}`;
      }
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
    this.refreshList();
  }

  async delete(vcdId: string) {
    await this.resourceService.removeItem(vcdId).then(() => {
      this.vo.alertText = `Successfully delete Resource with vcdId ${vcdId}`;
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
    this.refreshList();
  }

  async resetAlert(time?: number) {
    window.setTimeout(() => {
      this.vo.alertText = '';
    }, time || NOTIFICATION_EXIST_TIME);
  }

  vendorList: Object = {};
  list$: Observable<Item[]> = of([]);
  opened = false;
  refreshInterval: number;
  selected: Observable<Item[]> = of([])

  async save() {
    await this.refreshList();
  }

  cancel() {
    this.opened = false;
  }

  async ngOnInit() {
    await this.refreshList();
  }

  async refreshList() {
    this.list$ = of(await this.resourceService.getList());
    this.refreshInterval = window.setInterval(async () => {
      this.list$ = of(await this.resourceService.getList());
    }, RESOURCE_USAGE_REFRESH_PERIOD);
    this.vendorList = Object(await this.resourceService.getVendorList())
  }

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

}
