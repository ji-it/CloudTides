import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';


@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    ClarityModule,
  ],
  exports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    ClarityModule,
  ],
  declarations: [],
})
export class BaseModule {

  static forRoot (): ModuleWithProviders {
    return {
      ngModule: BaseModule,
      providers: [],
    };
  }

  static forChild (): ModuleWithProviders {
    return {
      ngModule: BaseModule,
      providers: [],
    };
  }
}

