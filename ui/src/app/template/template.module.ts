import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { TemplateRoutingModule, declarations, providers } from './template-routing.module';

@NgModule({
  declarations: [
    ...declarations,
  ],
  providers: [
    ...providers,
  ],
  imports: [
    SharedModule,
    TemplateRoutingModule,
  ],
})
export class TemplateModule {}
