import { Component, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';

import { TemplateService, Item } from '../template.service';

@Component({
  selector: 'tide-template-list',
  templateUrl: './template-list.component.html',
  styleUrls: ['./template-list.component.scss'],
})
export class TemplateListComponent implements OnInit {

  constructor(
    private readonly templateService: TemplateService,
  ) { }

  list$: Observable<Item[]> = of([]);

  ngOnInit() {
    this.list$ = this.templateService.getList();
  }

  add() {

  }

  onDelete(app: Item) {

  }

  onEdit(app: Item) {

  }
}
