import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { base } from '@tide-environments/base'

@Injectable()
export class TemplateService {

  constructor(
    private readonly http: HttpClient,
  ) {
  }

  private prefix = `${base.apiPrefix}/templates`;

  getList() {
    return this.http.get<Item[]>(`${this.prefix}`).pipe(
      map(mapList),
    );
  }

  addItem(payload: ItemPayload) {
    return this.http.post<ItemDTO>(`${this.prefix}`, payload).pipe(
      map(mapItem),
    );
  }

  editItem(id: string, payload: ItemPayload) {
    return this.http.put<ItemDTO>(`${this.prefix}/${id}`, payload).pipe(
      map(mapItem),
    );
  }

  removeItem(id: string) {
    return this.http.delete<ItemDTO>(`${this.prefix}/${id}`);
  }
}

// Raw
interface ItemDTO {
  name: string;
  displayName: string;
  version: string;
  date: string;
  url: string;
  cpu: number;
  mem: number;
  disk: number;
}

function mapList(raw: ItemDTO[]): Item[] {
  return raw.map(mapItem);
}

function mapItem(raw: ItemDTO): Item {
  return raw;
}

// UI
export type Item = ItemDTO;

export interface ItemPayload {
  name: string;
  displayName: string;
}

