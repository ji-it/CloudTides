import { ModuleWithProviders, NgModule } from '@angular/core';

import { BaseModule } from './base.module';
import { sharedComponents } from './component';

@NgModule({
  imports: [
    BaseModule,
  ],
  exports: [
    BaseModule,
    ...sharedComponents,
  ],
  declarations: [
    ...sharedComponents,
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
