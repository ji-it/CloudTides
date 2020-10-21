import { Component, OnInit, Input } from '@angular/core';
import { Item } from '../resource.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'tide-resource-card',
  templateUrl: './resource-card.component.html',
  styleUrls: ['./resource-card.component.scss'],
})
export class ResourceCardComponent implements OnInit {

  constructor(
    public readonly translate: TranslateService,
  ) { }

  @Input() item: Item;

  ngOnInit() {

  }

}
