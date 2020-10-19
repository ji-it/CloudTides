import { ModuleWithProviders, NgModule } from '@angular/core';

import { BaseModule } from './base.module';
import { sharedComponents } from './component';
import { sharedPipes } from './pipe';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { HttpClient } from '@angular/common/http';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';

// AoT requires an exported function for factories
export function HttpLoaderFactory(http: HttpClient) {
  return new TranslateHttpLoader(http);
}

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
