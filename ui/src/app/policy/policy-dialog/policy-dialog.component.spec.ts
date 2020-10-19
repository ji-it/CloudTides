import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PolicyDialogComponent } from './policy-dialog.component';

describe('policyDialogComponent', () => {
  let component: PolicyDialogComponent;
  let fixture: ComponentFixture<PolicyDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PolicyDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PolicyDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
