import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { Observable, of } from 'rxjs';
import { TemplateListComponent } from '../template-list/template-list.component';
import { TemplateService, ItemPayloadVM} from '../template.service';
import { VMCardComponent } from '../vm-card/vm-card.component';

@Component({
  selector: 'tide-template-dialogupdate',
  templateUrl: './template-dialogupdate.component.html',
  styleUrls: ['./template-dialogupdate.component.scss']
})
export class TemplateDialogUpdateComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly templateService: TemplateService,
    private readonly  vmCard: VMCardComponent,
  ) {
    this.templateForm = this.fb.group({
      id: [],
      disk: [vmCard.updateCPU, Validators.required],
      vmem: [vmCard.updateMem, Validators.required],
      vcpu: [vmCard.updateDisk, Validators.required],
      ports: [''],
      //templateID: ['', Validators.required],
    })
  }

  @Input() opened = false;
  @Input() vmid: number
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
    const payload = this.changeID(value)
    this.resetModal();
    this.vo.spinning = true;
    await this.templateService.editItemVM(payload).then(() => {
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

  private changeID (payload: ItemPayload) {
    const result : ItemPayload = {
      id: this.vmid,
      disk: payload.disk,
      vmem: payload.vmem,
      vcpu: payload.vcpu,
      ports: payload.ports,
    }
    return result;
  }

}

interface ItemPayload {
  id: number,
  disk: number,
  vmem: number,
  vcpu: number,
  ports: string,
}