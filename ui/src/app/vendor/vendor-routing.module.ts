import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { VendorComponent } from './vendor.component';
import { VendorService } from './vendor.service';
import { VendorListComponent } from './vendor-list/vendor-list.component';
//import { TemplateCardComponent } from './template-card/template-card.component';

const routes: Routes = [
  {
    path: '',
    component: VendorComponent,
    children: [
      {
        path: '',
        component: VendorListComponent,
      },
    ],
  },
];

export const declarations = [
  VendorComponent,
  VendorListComponent,
];

export const providers = [
  VendorService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class VendorRoutingModule {}
