import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { Observable, of } from 'rxjs';
import { TemplateListComponent } from '../template-list/template-list.component';
import { TemplateService } from '../template.service';

@Component({
  selector: 'tide-template-dialogvm',
  templateUrl: './template-dialogvm.component.html',
  styleUrls: ['./template-dialogvm.component.scss']
})
export class TemplateDialogVMComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly templateService: TemplateService,
    private readonly  templateList: TemplateListComponent,
  ) {
    this.tempList = Object.keys(templateList.TemplateList);
    this.template = templateList.TemplateList;
    this.templateForm = this.fb.group({
      name: ['', Validators.required],
      disk: ['', Validators.required],
      vmem: ['', Validators.required],
      vcpu: ['', Validators.required],
      templateID: ['', Validators.required],
    })
  }

  @Input() opened = false;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  templateForm: FormGroup;
  tempList: string[];
  template: any;

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
    const { value } = this.templateForm;
    this.resetModal();
    this.vo.spinning = true;
    await this.templateService.addItemVM(value).then(() => {
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