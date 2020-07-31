import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, pluck } from 'rxjs/operators';
import { base } from '@tide-environments/base'

@Injectable()
export class ResourceService {

  constructor(
    private readonly http: HttpClient,
  ) {
  }

  private prefix = `${base.apiPrefix}/resource/list`;

  getList() {
    return this.http.get<Item[]>(`${this.prefix}`).pipe(
      pluck('results'),
      map(mapList),
    );
  }

  addItem(payload: ItemPayload) {
    return this.http.post<ItemDTO>(`${this.prefix}`, `name=${payload.name}`).pipe(
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
  id: string;
  name: string;
  currentCPU: number;
  currentRAM: number;
  CPUPercent: number;
  RAMPercent: number;
}

function mapList(raw: ItemDTO[]): Item[] {
  return raw.map(mapItem);
}

function mapItem(raw: ItemDTO): Item {
  return raw;
}

// UI
export interface ItemPayload {
  name: string;
}

export type Item = ItemDTO;
