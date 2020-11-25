import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

import { ItemPayload } from '../resource.service';
import { TranslateService } from '@ngx-translate/core';
import { cloudPlatform, defaultCloudPlatformURL } from '@tide-config/cloudPlatform';

@Component({
  selector: 'tide-resource-dialog',
  templateUrl: './resource-dialog.component.html',
  styleUrls: ['./resource-dialog.component.scss'],
})
export class ResourceDialogComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
  ) {
    this.resourceForm = this.fb.group({
      href: [defaultCloudPlatformURL, [Validators.required]],
      datacenter: [''],
      org: [''],
      username: ['', Validators.required],
      password: ['', Validators.required],
    });
    this.cloudPlatformList = Object.keys(cloudPlatform);
    this.cloudPlatform = cloudPlatform;
  }

  @Input() opened = false;
  @Output() save = new EventEmitter<ItemPayload>();
  @Output() cancel = new EventEmitter();

  resourceForm: FormGroup;
  cloudPlatformList: string[];
  cloudPlatform: any;

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
