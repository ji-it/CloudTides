import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable()
export class PolicyService {

  constructor(
    private readonly http: HttpClient,
  ) {
  }

  private prefix = '/api/policies';

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
  id: string;
  name: string;
  project: string;
  type: string;
  cpu: number;
  memory: number;
  destroy: boolean;
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
