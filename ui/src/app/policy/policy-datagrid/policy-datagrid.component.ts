import { Component, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { Item, PolicyService } from '../policy.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'policy-policy-datagrid',
  templateUrl: './policy-datagrid.component.html',
  styleUrls: ['./policy-datagrid.component.scss'],
})
export class PolicyDatagridComponent implements OnInit {

  constructor(
    private policyService: PolicyService,
    public readonly translate: TranslateService,
  ) { }

  readonly vo = {
    pageSize: 20,
  };

  list$: Observable<Item[]> = of([]);
  opened = false;

  add(resource: Item) {
    this.policyService.addItem(resource).subscribe(item => {
      this.refreshList();
    });
  }

  cancel() {
    this.opened = false;
  }

  edit(item: Item) {

  }

  delete(item: Item) {

  }

  ngOnInit() {
    this.refreshList();
  }

  refreshList() {
    this.list$ = this.policyService.getList();
  }

}
