import { NgModule } from '@angular/core';
import { ChartsModule } from 'ng2-charts';

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
    HomeRoutingModule,
    ChartsModule,
  ],
})
export class HomeModule {}
