import { Component, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { Item, ResourceService } from '../resource.service';


@Component({
  selector: 'tide-resource-list',
  templateUrl: './resource-list.component.html',
  styleUrls: ['./resource-list.component.scss']
})
export class ResourceListComponent implements OnInit {

  constructor(
    private resourceService: ResourceService,
  ) {

  }

  list$: Observable<Item[]> = of([]);

  add () {
    this.resourceService.addItem({
      name: '',
      description: '',
    }).subscribe(item => {
      this.refreshList();
    });
  }

  ngOnInit() {
    this.refreshList();
  }

  refreshList() {
    this.list$ = this.resourceService.getList();
  }

}
