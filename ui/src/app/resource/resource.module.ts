import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { ResourceRoutingModule, declarations, providers } from './resource-routing.module';

@NgModule({
  declarations: [
    ...declarations,
  ],
  providers: [
    ...providers,
  ],
  imports: [
    SharedModule,
    ResourceRoutingModule,
  ],
})
export class ResourceModule {}
