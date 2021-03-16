import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { VendorRoutingModule, declarations, providers } from './vendor-routing.module';

@NgModule({
  declarations: [
    ...declarations,
  ],
  providers: [
    ...providers,
  ],
  imports: [
    SharedModule,
    VendorRoutingModule,
  ],
})
export class VendorModule {}
