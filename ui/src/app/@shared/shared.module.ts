import { ModuleWithProviders, NgModule } from '@angular/core';

import { BaseModule } from './base.module';
import { sharedComponents } from './component';
import { sharedPipes } from './pipe';
import { TranslateModule } from '@ngx-translate/core';

@NgModule({
  imports: [
    BaseModule,
  ],
  exports: [
    BaseModule,
    TranslateModule,
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
