import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
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
export class VMCardComponent implements OnInit {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly templateService: TemplateService,
    private readonly  templateList: TemplateListComponent,
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
    this.list$ = of(await this.templateService.getVMList(this.templateid));
    this.refreshInterval = window.setInterval(async () => {
      this.list$ = of(await this.templateService.getVMList(this.templateid));
    }, VENDOR_USAGE_REFRESH_PERIOD);
  }


  @Input() opened = false;
  @Input() templateid = 1;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  list$: Observable<ItemVM[]> = of([]);
  refreshInterval: number;

  ngOnInit(): void {
  }

  onCancel() {
    this.close();
  }

  private close() {
    this.cancel.emit();
  }

  delete(id: number){}

}