import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'cp-spinner',
  templateUrl: './spinner.component.html',
  styleUrls: ['./spinner.component.scss'],
})
export class SpinnerComponent implements OnInit {

  constructor() { }

  @Input()
  type: SpinnerType = 'global';

  ngOnInit() { }

}

export type SpinnerType = 'block' | 'global' | 'inline';
