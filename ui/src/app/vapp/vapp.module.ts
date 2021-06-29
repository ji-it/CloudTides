import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { VappRoutingModule, declarations, providers } from './vapp-routing.module';
import { VappListComponent } from './vapp-list/vapp-list.component';

@NgModule({
  declarations: [
    ...declarations,
    VappListComponent,
  ],
  providers: [
    ...providers,
  ],
  imports: [
    SharedModule,
    VappRoutingModule,
  ],
})
export class VappModule {}
