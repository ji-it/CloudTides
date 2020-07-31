import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { SharedModule } from '@tide-shared/shared.module';
import { interceptorProviders } from '@tide-shared/interceptor/api.interceptor';

import { AppRoutingModule, declarations, providers } from './app-routing.module';
import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent,
    ...declarations,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    SharedModule,
    AppRoutingModule
  ],
  providers: [
    ...providers,
    ...interceptorProviders,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
