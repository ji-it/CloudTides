import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core'
;
import { SharedModule } from '@tide-shared/shared.module';

import { AppRoutingModule, declarations, providers } from './app-routing.module';
import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent,
    ...declarations,
  ],
  imports: [
    BrowserModule,
    SharedModule,
    AppRoutingModule
  ],
  providers: [
    ...providers,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
