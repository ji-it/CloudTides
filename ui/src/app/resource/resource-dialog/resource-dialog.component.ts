import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

import { ItemPayload } from '../resource.service';

@Component({
  selector: 'tide-resource-dialog',
  templateUrl: './resource-dialog.component.html',
  styleUrls: ['./resource-dialog.component.scss']
})
export class ResourceDialogComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
  ) {
    this.resourceForm = this.fb.group({
      name: [ '', Validators.required ],
      datacenter: [ '' ],
      cluster: [ '' ],
      username: [ '', Validators.required ],
      password: [ '', Validators.required ],
    });
  }

  @Input() opened = false;
  @Output() save = new EventEmitter<ItemPayload>();
  @Output() cancel = new EventEmitter();

  resourceForm: FormGroup

  ngOnInit(): void {

  }

  onCancel() {
    this.close();
  }

  onSave() {
    this.save.emit(this.resourceForm.value);
    this.close();
  }

  private close() {
    this.cancel.emit();
  }

}
