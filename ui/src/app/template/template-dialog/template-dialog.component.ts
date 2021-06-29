import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { Observable, of } from 'rxjs';
import { TemplateListComponent } from '../template-list/template-list.component';
import { TemplateService } from '../template.service';

@Component({
  selector: 'tide-template-dialog',
  templateUrl: './template-dialog.component.html',
  styleUrls: ['./template-dialog.component.scss']
})
export class TemplateDialogComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly templateService: TemplateService,
    private readonly  templateList: TemplateListComponent,
  ) {
    //this.resourceList = Object.keys(templateList.ResList);
    this.templateForm = this.fb.group({
      type: ['', Validators.required],
      name: ['', Validators.required],
      tag: ['', Validators.required],
      description: ['', Validators.required],
      path: ['', Validators.required]
    })
  
  }

  @Input() opened = false;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  templateForm: FormGroup;
  //resourceList: string[];
  resource: any;
  resourceList: string[] = [`Big Data`, `Machine Learning`, `Research`, `Cloud Native`];
  typeList: string[] = [`OVA`, `Others`]

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
    await this.templateService.addItem(value).then(() => {
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