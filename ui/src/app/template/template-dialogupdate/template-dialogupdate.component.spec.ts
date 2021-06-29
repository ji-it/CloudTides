import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TemplateDialogUpdateComponent } from './template-dialogupdate.component';

describe('TemplateDialogUpdateComponent', () => {
  let component: TemplateDialogUpdateComponent;
  let fixture: ComponentFixture<TemplateDialogUpdateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TemplateDialogUpdateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TemplateDialogUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
