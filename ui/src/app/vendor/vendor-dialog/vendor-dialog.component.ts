import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { defaultType, vendorTypes } from '@tide-shared/config/vendorType';
import { VendorService } from '../vendor.service';

@Component({
  selector: 'tide-vendor-dialog',
  templateUrl: './vendor-dialog.component.html',
  styleUrls: ['./vendor-dialog.component.scss']
})
export class VendorDialogComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    public readonly vendorService: VendorService,
  ) {
    this.vendorForm = this.fb.group({
      name: ['', Validators.required],
      url: ['', Validators.required],
      vendorType: [defaultType, Validators.required],
      version: [''],
    });
    this.vendorTypeList = Object.keys(vendorTypes);
    this.vendorTypes = vendorTypes;
  }

  @Input() opened = false;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  vendorForm: FormGroup;
  vendorTypeList: string[];
  vendorTypes: any;

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
    const { value } = this.vendorForm;
    this.resetModal();
    this.vo.spinning = true;
    await this.vendorService.addItem(value).then(() => {
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
