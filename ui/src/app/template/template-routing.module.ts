import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { TemplateComponent } from './template.component';
import { TemplateService } from './template.service';
import { TemplateListComponent } from './template-list/template-list.component';
import { TemplateCardComponent } from './template-card/template-card.component';

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
  TemplateCardComponent,
];

export const providers = [
  TemplateService,
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class TemplateRoutingModule {}
