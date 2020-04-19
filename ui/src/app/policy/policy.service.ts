import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { base } from '@tide-environments/base'

@Injectable()
export class PolicyService {

  constructor(
    private readonly http: HttpClient,
  ) {
  }

  private prefix = `${base.apiPrefix}/schedulerPolicie`;

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

// DTO
interface ItemDTO {
  crid: string;
  template: string;
  startConditions: (string | number)[][];
  stopConditions: (string | number)[][];
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

}
