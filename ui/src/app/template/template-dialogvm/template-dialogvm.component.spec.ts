import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TemplateDialogVMComponent } from './template-dialogvm.component';

describe('TemplateDialogVMComponent', () => {
  let component: TemplateDialogVMComponent;
  let fixture: ComponentFixture<TemplateDialogVMComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TemplateDialogVMComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TemplateDialogVMComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
