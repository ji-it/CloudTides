import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { Item, ItemPayload, ResourceService } from '../resource.service';
import { TranslateService } from '@ngx-translate/core';
import { RESOURCE_USAGE_REFRESH_PERIOD } from '@tide-config/const';

@Component({
  selector: 'tide-resource-list',
  templateUrl: './resource-list.component.html',
  styleUrls: ['./resource-list.component.scss'],
})
export class ResourceListComponent implements OnInit, OnDestroy {

  constructor(
    private resourceService: ResourceService,
    public readonly translate: TranslateService,
  ) {

  }

  list$: Observable<Item[]> = of([]);
  opened = false;
  refreshInterval: number;

  add(resource: ItemPayload) {
    this.resourceService.addItem(resource).subscribe(async item => {
      await this.refreshList();
    });
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
  }

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

}
