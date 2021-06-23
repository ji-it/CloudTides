import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { TemplateComponent } from './template.component';
import { TemplateService } from './template.service';
import { TemplateListComponent } from './template-list/template-list.component';
import { TemplateDialogComponent } from './template-dialog/template-dialog.component';
import { TemplateDialogVMComponent } from './template-dialogvm/template-dialogvm.component';
import { VMCardComponent } from './vm-card/vm-card.component';
import { TemplateDialogUpdateComponent } from './template-dialogupdate/template-dialogupdate.component';

const routes: Routes = [
  {
    path: '',
    component: TemplateComponent,
    children: [
      {
        path: '',
        component: TemplateListComponent,
      },
    ],
  },
];

export const declarations = [
  TemplateComponent,
  TemplateListComponent,
  TemplateDialogComponent,
  TemplateDialogVMComponent,
  VMCardComponent,
  TemplateDialogUpdateComponent,
];

export const providers = [
  TemplateService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class TemplateRoutingModule {}
