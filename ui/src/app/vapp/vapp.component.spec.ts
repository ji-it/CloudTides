import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VappComponent } from './vapp.component';

describe('VappComponent', () => {
  let component: VappComponent;
  let fixture: ComponentFixture<VappComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VappComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VappComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
