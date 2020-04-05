import { Component, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { Item, PolicyService } from '../policy.service';

@Component({
  selector: 'policy-policy-datagrid',
  templateUrl: './policy-datagrid.component.html',
  styleUrls: ['./policy-datagrid.component.scss']
})
export class PolicyDatagridComponent implements OnInit {

  constructor(
    private policyService: PolicyService,
  ) { }

  readonly vo = {
    pageSize: 20,
  };

  list$: Observable<Item[]> = of([]);

  add() {

  }

  edit(item: Item) {

  }

  delete(item: Item) {

  }

  ngOnInit() {
    this.list$ = this.policyService.getList();
  }

}
