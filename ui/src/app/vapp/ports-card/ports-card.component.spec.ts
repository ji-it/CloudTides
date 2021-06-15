import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PortsCardComponent } from './ports-card.component';

describe('TemplateDialogVMComponent', () => {
  let component: PortsCardComponent;
  let fixture: ComponentFixture<PortsCardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PortsCardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PortsCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
