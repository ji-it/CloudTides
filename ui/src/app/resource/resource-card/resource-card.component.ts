import { Component, OnInit, Input } from '@angular/core';
import { Item, ResourceService } from '../resource.service';
import { TranslateService } from '@ngx-translate/core';
import { NOTIFICATION_EXIST_TIME } from '@tide-config/const';
import { LoginService } from '../../login/login.service';

@Component({
  selector: 'tide-resource-card',
  templateUrl: './resource-card.component.html',
  styleUrls: ['./resource-card.component.scss'],
})
export class ResourceCardComponent implements OnInit {

  constructor(
    public readonly translate: TranslateService,
    public readonly resourceService: ResourceService,
    public readonly loginService: LoginService,
  ) { }

  @Input() item: Item;

  readonly vo = {
    alertType: '',
    alertText: '',
  };

  ngOnInit() {

  }

  async contribute() {
    await this.resourceService.contributeResource(this.item.id).then((resp) => {
      if (resp.contributed) {
        this.vo.alertText = `Successfully start contributing Resource ${this.item.name}`;
      } else {
        this.vo.alertText = `Successfully stop contributing Resource ${this.item.name}`;
      }
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
  }

  async activate() {
    await this.resourceService.activateResource(this.item.id).then((resp) => {
      if (resp.activated) {
        this.vo.alertText = `Successfully activate Resource ${this.item.name}`;
      } else {
        this.vo.alertText = `Successfully deactivate Resource ${this.item.name}`;
      }
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
  }

  async delete() {
    await this.resourceService.removeItem(this.item.vcdId).then(() => {
      this.vo.alertText = `Successfully delete Resource ${this.item.name}`;
      this.vo.alertType = 'success';
    }, (error) => {
      this.vo.alertType = 'danger';
      this.vo.alertText = error;
    }).then(() => {
      this.resetAlert();
    });
  }

  async resetAlert(time?: number) {
    window.setTimeout(() => {
      this.vo.alertText = '';
    }, time || NOTIFICATION_EXIST_TIME);
  }

}
