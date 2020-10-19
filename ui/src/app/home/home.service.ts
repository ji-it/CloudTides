import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class HomeService {

  constructor(
    private readonly http: HttpClient,
  ) {
  }

}
