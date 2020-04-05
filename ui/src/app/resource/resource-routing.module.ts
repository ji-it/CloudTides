import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ResourceComponent } from './resource.component';
import { ResourceService } from './resource.service';
import { ResourceListComponent } from './resource-list/resource-list.component';
import { ResourceCardComponent } from './resource-card/resource-card.component';


const routes: Routes = [
  {
    path: '',
    component: ResourceComponent,
    children: [
      {
        path: '',
        component: ResourceListComponent,
      }
    ]
  }
];

export const declarations = [
  ResourceComponent,
  ResourceListComponent,
  ResourceCardComponent,
];

export const providers = [
  ResourceService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ResourceRoutingModule { }
