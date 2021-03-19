import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { Item, ItemPayload, VendorService } from '../vendor.service';
import { TranslateService } from '@ngx-translate/core';
import { VENDOR_USAGE_REFRESH_PERIOD } from '@tide-config/const';

@Component({
  selector: 'tide-vendor-list',
  templateUrl: './vendor-list.component.html',
  styleUrls: ['./vendor-list.component.scss'],
})
export class VendorListComponent implements OnInit, OnDestroy {

  constructor(
    private vendorService: VendorService,
    public readonly translate: TranslateService,
  ) {

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
    await this.refreshList();
  }

  async refreshList() {
    this.list$ = of(await this.vendorService.getList());
    this.refreshInterval = window.setInterval(async () => {
      this.list$ = of(await this.vendorService.getList());
    }, VENDOR_USAGE_REFRESH_PERIOD);
  }

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

}
