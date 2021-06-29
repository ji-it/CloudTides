import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { NOTIFICATION_EXIST_TIME, VENDOR_USAGE_REFRESH_PERIOD } from '@tide-shared/config/const';
import { Observable, of } from 'rxjs';
import { VappListComponent } from '../vapp-list/vapp-list.component';
import { ItemPort, VappService } from '../vapp.service';

@Component({
  selector: 'tide-vm-card',
  templateUrl: './vm-card.component.html',
  styleUrls: ['./vm-card.component.scss']
})
export class VMCardComponent implements OnInit, OnDestroy {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly vappService: VappService,
    public VappList: VappListComponent,
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
    this.VappList.VMlist$ = of(await this.vappService.getVMList(this.VappList.VappID));
    this.refreshInterval = window.setInterval(async () => {
      this.VappList.VMlist$ = of(await this.vappService.getVMList(this.VappList.VappID));
    }, VENDOR_USAGE_REFRESH_PERIOD);
  }


  @Input() opened = false;
  @Input() vappid = 1;
  @Output() save = new EventEmitter();
  @Output() cancel = new EventEmitter();

  refreshInterval: number;
  Portslist$: Observable<ItemPort[]> = of([]);
  VMID = 1;
  PortOpened = false;

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

  async ngOnInit() {
    await this.refreshList();
  }

  cancelPort() {
    this.PortOpened = false;
  }

  onCancel() {
    this.close();
  }

  private close() {
    this.cancel.emit();
  }

  async displayPort(id: number) {
    this.VMID = id;
    this.PortOpened = true;
    this.Portslist$ = of(await this.vappService.getPortList(id))
  }
  /*async delete(id: number) {
    await this.vappService.removeItemVM(id).then(() => {
      this.vo.alertText = `Successfully delete vendor with id ${id}`;
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
    this.refreshList();
  }*/

}