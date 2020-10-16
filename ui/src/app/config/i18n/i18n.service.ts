import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { i18nSupportList } from './i18nSupportList';

@Injectable()
export class I18nService {
  constructor(
    public readonly translate: TranslateService,
  ) {
    this.resetLanguage();
    translate.use(localStorage.getItem('i18n'));
  }

  public readonly i18nSupportList = i18nSupportList;

  i18nChoice: string;

  resetLanguage(): any {
    if (!localStorage.getItem('i18n')) {
      localStorage.setItem('i18n', 'en');
    }
    this.i18nChoice = localStorage.getItem('i18n');
  }

  getLanguage(): string {
    this.resetLanguage();
    return this.i18nChoice;
  }

  setLanguage(): any {
    localStorage.setItem('i18n', this.i18nChoice);
    this.translate.use(this.i18nChoice);
  }
}
