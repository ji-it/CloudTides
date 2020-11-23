import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { base } from '@tide-environments/base';
import { VCD_URL_PATH } from '@tide-config/path';
import { LoginService } from '../login/login.service';
import toFixed from 'accounting-js/lib/toFixed.js';

@Injectable()
export class ResourceService {

  constructor(
    private readonly http: HttpClient,
    private readonly loginService: LoginService,
  ) {
  }

  private prefix = `${base.apiPrefix}/computeResource`;

  async getList() {
    const list = await this.http.get<Item[]>(base.apiPrefix + VCD_URL_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const usage: Item[] = [];
    for (const resource of list) {
      const rawUsage = await this.http.get<ItemUsage>(`${base.apiPrefix}/usage/${resource.id}`, {
        headers: {
          Authorization: `Bearer ${this.loginService.token}`,
        },
      }).toPromise();
      const resourceItem: Item = {
        id: resource.id,
        name: rawUsage.name,
        cpu: rawUsage.totalCPU / 1000,
        mem: rawUsage.totalRAM / 1000,
        disk: NaN,
        usage: {
          'cpu%': toFixed(rawUsage.percentCPU * 100, 2),
          'mem%': toFixed(rawUsage.percentRAM * 100, 2),
          'disk%': toFixed(NaN, 2),
        },
      };
      usage.push(resourceItem);
    }
    return usage;
  }

  addItem(payload: ItemPayload) {
    const body = {
      ...payload,
      policy: 0,
    };
    return this.http.post<any>(base.apiPrefix + VCD_URL_PATH, body, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    });
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
interface ItemUsage {
  currentCPU: number;
  totalCPU: number;
  currentRAM: number;
  totalRAM: number;
  percentCPU: number;
  percentRAM: number;
  name: string;
}

interface ItemDTO {
  id: string;
  name: string;
  // unit: GHz
  cpu: number;
  // unit: GB
  mem: number;
  // unit: GB
  disk: number;
  usage: {
    'cpu%': number;
    'mem%': number;
    'disk%': number;
  }
}

function mapList(raw: ItemDTO[]): Item[] {
  return raw.map(mapItem);
}

function mapItem(raw: ItemDTO): Item {
  return raw;
}

// UI
export interface ItemPayload {
  datacenter: string;
  name: string;
  org: string;
  username: string,
  password: string,
}

export type Item = ItemDTO;
