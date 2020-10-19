import { Component, OnInit, Input } from '@angular/core';
import { Item } from '../resource.service';

@Component({
  selector: 'tide-resource-card',
  templateUrl: './resource-card.component.html',
  styleUrls: ['./resource-card.component.scss'],
})
export class ResourceCardComponent implements OnInit {

  constructor() { }

  @Input() item: Item;

  ngOnInit() {

  }

}
