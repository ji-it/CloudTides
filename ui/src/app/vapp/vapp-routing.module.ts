import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { VappComponent } from './vapp.component';
import { VappService } from './vapp.service';
import { VappListComponent } from './vapp-list/vapp-list.component';
import { VappDialogComponent } from './vapp-dialog/vapp-dialog.component';
import { VMCardComponent } from './vm-card/vm-card.component';
import { PortsCardComponent } from './ports-card/ports-card.component';
//import { TemplateCardComponent } from './template-card/template-card.component';

const routes: Routes = [
  {
    path: '',
    component: VappComponent,
    children: [
      {
        path: '',
        component: VappListComponent,
      },
    ],
  },
];

export const declarations = [
  VappComponent,
  VappListComponent,
  VappDialogComponent,
  VMCardComponent,
  PortsCardComponent,
];

export const providers = [
  VappService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class VappRoutingModule {}
