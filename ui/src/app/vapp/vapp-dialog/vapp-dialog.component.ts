import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { VappService } from '../vapp.service';

@Component({
  selector: 'tide-vapp-dialog',
  templateUrl: './vapp-dialog.component.html',
  styleUrls: ['./vapp-dialog.component.scss']
})
export class VappDialogComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly vappService: VappService,
  ) {
    this.vappForm = this.fb.group({
      name: ['', Validators.required],
      template: ['', Validators.required],
      vendor: ['', Validators.required],
      datacenter: ['', Validators.required],
    })
  }

  @Input() opened = false;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  vappForm: FormGroup;
  vappTypeList: string[];
  vappTypes: any;

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
    const { value } = this.vappForm;
    this.resetModal();
    this.vo.spinning = true;
    await this.vappService.addItem(value).then(() => {
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
