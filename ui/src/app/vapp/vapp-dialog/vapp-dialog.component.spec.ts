import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VappDialogComponent } from './vapp-dialog.component';

describe('VappDialogComponent', () => {
  let component: VappDialogComponent;
  let fixture: ComponentFixture<VappDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VappDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VappDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
