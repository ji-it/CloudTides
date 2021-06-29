import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { LoginComponent } from './login/login.component';
import { LoginService } from './login/login.service';
import { RegisterComponent } from './register/register.component';

import {
  LOGIN_PATH_NAME,
  HOME_PATH_NAME,
  RESOURCE_PATH_NAME,
  TEMPLATE_PATH_NAME,
  POLICY_PATH_NAME,
  REGISTER_PATH_NAME,
  VENDOR_PATH_NAME,
  VAPP_PATH_NAME,
} from '@tide-config/path';

import { AuthGuard } from '@tide-guard/auth.guard';
import { RegisterService } from './register/register.service';

const routes: Routes = [
  {
    path: '',
    canActivateChild: [AuthGuard],
    children: [
      {
        path: '',
        pathMatch: 'full',
        // redirectTo: HOME_PATH_NAME,
        redirectTo: RESOURCE_PATH_NAME
      },
      {
        path: LOGIN_PATH_NAME,
        component: LoginComponent,
        data: {
          anonymous: true,
        } as RouterData,
      },
      {
        path: REGISTER_PATH_NAME,
        component: RegisterComponent,
        data: {
          anonymous: true,
        } as RouterData,
      },
      {
        path: HOME_PATH_NAME,
        loadChildren: () => import('./home/home.module').then(m => m.HomeModule),
      },
      {
        path: VENDOR_PATH_NAME,
        loadChildren: () => import('./vendor/vendor.module').then(m => m.VendorModule)
      },
      {
        path: VAPP_PATH_NAME,
        loadChildren: () => import('./vapp/vapp.module').then(m => m.VappModule)
      },
      {
        path: RESOURCE_PATH_NAME,
        loadChildren: () => import('./resource/resource.module').then(m => m.ResourceModule),
      },
      {
        path: POLICY_PATH_NAME,
        loadChildren: () => import('./policy/policy.module').then(m => m.PolicyModule),
      },
      {
        path: TEMPLATE_PATH_NAME,
        loadChildren: () => import('./template/template.module').then(m => m.TemplateModule),
      },
    ],
  },
];

export const declarations = [
  LoginComponent,
  RegisterComponent,
];

export const providers = [
  AuthGuard,
  LoginService,
  RegisterService,
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}

export interface RouterData {
  anonymous?: boolean;
}
