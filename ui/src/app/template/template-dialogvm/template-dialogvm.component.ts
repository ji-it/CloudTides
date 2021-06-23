import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { Observable, of } from 'rxjs';
import { TemplateListComponent } from '../template-list/template-list.component';
import { TemplateService, ItemPayloadVM} from '../template.service';

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
    //this.resourceList = Object.keys(templateList.ResList);
    //this.resource = templateList.ResList;
    this.templateForm = this.fb.group({
      name: ['', Validators.required],
      disk: ['', Validators.required],
      vmem: ['', Validators.required],
      vcpu: ['', Validators.required],
      ports: [''],
      //templateID: ['', Validators.required],
    })
  }

  @Input() opened = false;
  @Input() templateid = 1;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  templateForm: FormGroup;
  tempList: string[];
  template: any;
  resourceList: string[];
  resource: any;

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
    const PayLoad = this.addTemplateID(value)
    this.resetModal();
    this.vo.spinning = true;
    await this.templateService.addItemVM(PayLoad).then(() => {
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

  private addTemplateID(payload: ItemPayload) {
    const result : ItemPayloadVM = {
      name: payload.name,
      disk: payload.disk,
      vmem: payload.vmem,
      vcpu: payload.vcpu,
      templateID: this.templateList.TemplateID,
      ports: payload.ports,
    }
    return result;
  }

}

interface ItemPayload {
  name: string,
  disk: number,
  vmem: number,
  vcpu: number,
  ports: string,
}