import { NgModule } from '@angular/core';

import { SharedModule } from '@tide-shared/shared.module';

import { TemplateRoutingModule, declarations, providers } from './template-routing.module';
import { TemplateDialogComponent } from './template-dialog/template-dialog.component';

@NgModule({
  declarations: [
    ...declarations,
    TemplateDialogComponent,
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
