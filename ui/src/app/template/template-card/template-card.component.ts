import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';

import { Item } from '../template.service';

@Component({
  selector: 'tide-template-card',
  templateUrl: './template-card.component.html',
  styleUrls: ['./template-card.component.scss'],
})
export class TemplateCardComponent implements OnInit {

  constructor() { }

  @Input() template: Item;
  @Output() delete = new EventEmitter<Item>();
  @Output() edit = new EventEmitter<Item>();

  onDelete(template: Item) {
    this.delete.emit(template);
  }

  onEdit(template: Item) {
    this.edit.emit(template);
  }

  ngOnInit() {
  }

}
