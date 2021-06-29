import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { NOTIFICATION_EXIST_TIME, VENDOR_USAGE_REFRESH_PERIOD } from '@tide-shared/config/const';
import { Observable, of } from 'rxjs';
import { VappListComponent } from '../vapp-list/vapp-list.component';
import { VappService } from '../vapp.service';
import { VMCardComponent } from '../vm-card/vm-card.component';

@Component({
  selector: 'tide-ports-card',
  templateUrl: './ports-card.component.html',
  styleUrls: ['./ports-card.component.scss']
})
export class PortsCardComponent implements OnInit, OnDestroy {

  constructor(
    private readonly fb: FormBuilder,
    public readonly translate: TranslateService,
    private readonly vappService: VappService,
    public VappList: VappListComponent,
    public VMCard: VMCardComponent,
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
    this.VMCard.Portslist$ = of(await this.vappService.getPortList(this.VMCard.VMID));
    this.refreshInterval = window.setInterval(async () => {
      this.VMCard.Portslist$ = of(await this.vappService.getPortList(this.VMCard.VMID));
    }, VENDOR_USAGE_REFRESH_PERIOD);
  }


  @Input() opened = false;
  @Input() vmid = 1;
  @Output() cancel = new EventEmitter();

  refreshInterval: number;

  ngOnDestroy(): void {
    window.clearInterval(this.refreshInterval);
  }

  async ngOnInit() {
    await this.refreshList();
  }

  onCancel() {
    this.close();
  }

  private close() {
    this.cancel.emit();
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