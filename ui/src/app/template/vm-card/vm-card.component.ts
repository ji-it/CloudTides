import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { NOTIFICATION_EXIST_TIME, VENDOR_USAGE_REFRESH_PERIOD } from '@tide-shared/config/const';
import { Observable, of } from 'rxjs';
import { TemplateListComponent } from '../template-list/template-list.component';
import { ItemVM, TemplateService } from '../template.service';

@Component({
  selector: 'tide-vm-card',
  templateUrl: './vm-card.component.html',
  styleUrls: ['./vm-card.component.scss']
})
export class VMCardComponent implements OnInit, OnDestroy {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly templateService: TemplateService,
    public TemplateList: TemplateListComponent,
  ) {
  }

  readonly vo = {
    alertType: '',
    alertText: '',
  };

  async resetAlert(time?: number) {
    window.setTimeout(() => {
      this.vo.alertText = '';
    }, time || NOTIFICATION_EXIST_TIME);
  }


  async refreshList() {
    this.TemplateList.VMlist$ = of(await this.templateService.getVMList(this.TemplateList.TemplateID));
    this.refreshInterval = window.setInterval(async () => {
      this.TemplateList.VMlist$ = of(await this.templateService.getVMList(this.TemplateList.TemplateID));
    }, VENDOR_USAGE_REFRESH_PERIOD);
  }


  @Input() opened = false;
  @Input() templateid = 1;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  list$: Observable<ItemVM[]> = of([]);
  refreshInterval: number;
  dialogOpened = false;
  updateId = 1;
  updateCPU = 1;
  updateMem = 1;
  updateDisk = 1;

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

  async ngOnInit() {
    await this.refreshList();
  }

  Cancel() {
    this.dialogOpened = false;
  }

  async Save() {
    await this.refreshList();
  }

  Click(id: number, cpu: number, mem: number, disk: number) {
    this.updateId = id;
    this.updateCPU = cpu;
    this.updateMem = mem;
    this.updateDisk = disk;
    this.dialogOpened = true;
  }

  onCancel() {
    this.close();
  }

  private close() {
    this.cancel.emit();
  }

  async delete(id: number) {
    await this.templateService.removeItemVM(id).then(() => {
      this.vo.alertText = `Successfully delete vendor with id ${id}`;
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
    this.refreshList();
  }

}