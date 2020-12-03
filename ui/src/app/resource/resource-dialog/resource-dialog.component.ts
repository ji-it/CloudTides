import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';

import { ItemPayload, ResourceService } from '../resource.service';
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
    public readonly resourceService: ResourceService,
  ) {
    this.resourceForm = this.fb.group({
      href: [defaultCloudPlatformURL, [Validators.required]],
      datacenter: [''],
      org: [''],
      network: [''],
      catalog: [''],
      username: ['', Validators.required],
      password: ['', Validators.required],
    });
    this.cloudPlatformList = Object.keys(cloudPlatform);
    this.cloudPlatform = cloudPlatform;
  }

  @Input() opened = false;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  resourceForm: FormGroup;
  cloudPlatformList: string[];
  cloudPlatform: any;

  readonly vo = {
    serverError: '',
    spinning: false,
  };

  ngOnInit(): void {

  }

  onCancel() {
    this.close();
  }

  async onSave() {
    const { value } = this.resourceForm;
    this.resetModal();
    this.vo.spinning = true;
    await this.resourceService.addItem(value).then(() => {
      this.save.emit('');
      this.close();
      this.vo.spinning = false;
    }, (error) => {
      this.vo.serverError = error;
      this.vo.spinning = false;
    });
  }

  private close() {
    this.cancel.emit();
  }

  private resetModal() {
    this.vo.serverError = '';
    this.vo.spinning = false;
  }

}
