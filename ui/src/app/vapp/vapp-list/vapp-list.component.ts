import { Component, OnInit } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { NOTIFICATION_EXIST_TIME, VENDOR_USAGE_REFRESH_PERIOD } from '@tide-shared/config/const';
import { Observable, of } from 'rxjs';
import { LoginService } from 'src/app/login/login.service';
import { Item, ItemPort, ItemResource, ItemVM, VappService } from '../vapp.service';

@Component({
  selector: 'tide-vapp-list',
  templateUrl: './vapp-list.component.html',
  styleUrls: ['./vapp-list.component.scss']
})
export class VappListComponent implements OnInit {

  constructor(
    public vappService: VappService,
    public readonly translate: TranslateService,
    public readonly loginService: LoginService,
  ) { }

  readonly vo = {
    alertType: '',
    alertText: '',
  };

  async delete(id: string) {
    await this.vappService.removeItem(id).then(() => {
      this.vo.alertText = `Successfully delete vapp with id ${id}`;
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

  VMlist$: Observable<ItemVM[]> = of([]);
  Portlist$: Observable<ItemPort[]> = of([]);
  VappID = 1;
  resourceList: Object = {};
  vendorList: Object = {};
  templateList: Object = {};
  list$: Observable<Item[]> = of([]);
  opened = false;
  displayOpened = false;
  refreshInterval: number;

  async save() {
    await this.refreshList();
  }

  cancel() {
    this.opened = false;
    this.displayOpened = false;
  }

  async ngOnInit() {
    await this.refreshList();
  }

  async displayVM(id: number) {
    this.VappID = id;
    this.displayOpened = true;
    this.VMlist$ = of(await this.vappService.getVMList(id))
  }

  async refreshList() {
    this.list$ = of(await this.vappService.getList());
    this.refreshInterval = window.setInterval(async () => {
      this.list$ = of(await this.vappService.getList());
    }, VENDOR_USAGE_REFRESH_PERIOD);
    this.vendorList = Object(await this.vappService.getVendorList());
    this.resourceList = Object(await this.vappService.getResourceList());
    this.templateList = Object(await this.vappService.getTemplateList());
  }

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }
}
