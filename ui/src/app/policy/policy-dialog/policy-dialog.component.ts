import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

import { ItemPayload } from '../policy.service';

@Component({
  selector: 'tide-policy-dialog',
  templateUrl: './policy-dialog.component.html',
  styleUrls: ['./policy-dialog.component.scss']
})
export class PolicyDialogComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
  ) {
    this.policyForm = this.fb.group({
      template: [ '', Validators.required ],
      startConditions: [ '', Validators.required ],
      stopConditions: [ '', Validators.required ],
    });
  }

  @Input() opened = false;
  @Output() save = new EventEmitter<ItemPayload>();
  @Output() cancel = new EventEmitter();

  policyForm: FormGroup

  ngOnInit(): void {

  }

  onCancel() {
    this.close();
  }

  onSave() {
    this.save.emit(this.policyForm.value);
    this.close();
  }

  private close() {
    this.cancel.emit();
  }

}
