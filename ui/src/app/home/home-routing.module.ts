import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './home.component';
import { HomeService } from './home.service';

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
    children: [],
  },
];

export const declarations = [
  HomeComponent,
];

export const providers = [
  HomeService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class HomeRoutingModule {}
