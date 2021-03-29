import { TestBed } from '@angular/core/testing';

import { VappService } from './vapp.service';

describe('VappService', () => {
  let service: VappService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(VappService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});