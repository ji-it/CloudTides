import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { VendorRoutingModule, declarations, providers } from './vendor-routing.module';
import { VendorDialogComponent } from './vendor-dialog/vendor-dialog.component';

@NgModule({
  declarations: [
    ...declarations,
    VendorDialogComponent,
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
