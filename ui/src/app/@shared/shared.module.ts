import { ModuleWithProviders, NgModule } from '@angular/core';

import { BaseModule } from './base.module';
import { sharedComponents } from './component';
import { sharedPipes } from './pipe';

@NgModule({
  imports: [
    BaseModule,
  ],
  exports: [
    BaseModule,
    ...sharedComponents,
    ...sharedPipes,
  ],
  declarations: [
    ...sharedComponents,
    ...sharedPipes,
  ],
})
export class SharedModule {

  static forRoot(): ModuleWithProviders {
    return {
      ngModule: SharedModule,
      providers: [],
    };
  }

  static forChild(): ModuleWithProviders {
    return {
      ngModule: SharedModule,
      providers: [],
    };
  }
}
