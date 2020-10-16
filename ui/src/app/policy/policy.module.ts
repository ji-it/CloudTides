import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { PolicyRoutingModule, declarations, providers } from './policy-routing.module';

@NgModule({
  declarations: [
    ...declarations,
  ],
  providers: [
    ...providers,
  ],
  imports: [
    SharedModule,
    PolicyRoutingModule,
  ],
})
export class PolicyModule {}
