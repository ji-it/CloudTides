import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VappListComponent } from './vapp-list.component';

describe('VappListComponent', () => {
  let component: VappListComponent;
  let fixture: ComponentFixture<VappListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VappListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VappListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
