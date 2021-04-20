import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VMCardComponent } from './vm-card.component';

describe('TemplateDialogVMComponent', () => {
  let component: VMCardComponent;
  let fixture: ComponentFixture<VMCardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VMCardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VMCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
