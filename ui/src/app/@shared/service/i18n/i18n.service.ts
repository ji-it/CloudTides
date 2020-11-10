import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { i18nSupportList } from '@tide-config/i18nSupportList';
import { LOCAL_STORAGE_KEY } from '@tide-config/const';

@Injectable()
export class I18nService {
  constructor(
    public readonly translate: TranslateService,
  ) {
    this.resetLanguage();
    translate.use(localStorage.getItem(LOCAL_STORAGE_KEY.I18N));
  }

  public readonly i18nSupportList = i18nSupportList;

  i18nChoice: string;

  resetLanguage(): any {
    if (!localStorage.getItem(LOCAL_STORAGE_KEY.I18N)) {
      localStorage.setItem(LOCAL_STORAGE_KEY.I18N, 'en');
    }
    this.i18nChoice = localStorage.getItem(LOCAL_STORAGE_KEY.I18N);
  }

  getLanguage(): string {
    this.resetLanguage();
    return this.i18nChoice;
  }

  setLanguage(): any {
    localStorage.setItem(LOCAL_STORAGE_KEY.I18N, this.i18nChoice);
    this.translate.use(this.i18nChoice);
  }
}
