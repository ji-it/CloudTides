import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { PolicyComponent } from './policy.component';
import { PolicyService } from './policy.service';
import { PolicyDatagridComponent } from './policy-datagrid/policy-datagrid.component'

const routes: Routes = [
  {
    path: '',
    component: PolicyComponent,
    children: [
      {
        path: '',
        component: PolicyDatagridComponent,
      }
    ]
  }
];

export const declarations = [
  PolicyComponent,
  PolicyDatagridComponent,
];

export const providers = [
  PolicyService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class PolicyRoutingModule { }
