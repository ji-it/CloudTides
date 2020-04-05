import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { HomeRoutingModule, declarations, providers } from './home-routing.module';


@NgModule({
  declarations: [
    ...declarations,
  ],
  providers: [
    ...providers,
  ],
  imports: [
    SharedModule,
    HomeRoutingModule
  ]
})
export class HomeModule { }
