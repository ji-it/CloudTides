import { Application } from 'express';

import * as fromAuth from './auth';
import * as fromResource from './resource';

export const runInterceptors = (server: Application ) => {
  server.post('/session', fromAuth.postLogin)
  server.post('/resources', fromResource.postHandler)
}
