import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from '@tide-environments/environment';
import { LoginService } from '../login/login.service';
import { TEMPLATE_PATH, VCD_URL_PATH } from '@tide-shared/config/path';

@Injectable()
export class TemplateService {

  constructor(
    private readonly http: HttpClient,
    private readonly loginService: LoginService,
  ) {
  }

  private prefix = `${environment.apiPrefix}/template`;

  async getList() {
    const TemplateList = await this.http.get<ItemDTO[]>(environment.apiPrefix + TEMPLATE_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const List: ItemDTO[] = [];
    for (const tem of TemplateList) {
      /*const res = await this.http.get<ItemRes>(environment.apiPrefix + VCD_URL_PATH + `/1`, {
        headers: {
          Authorization: `Bearer ${this.loginService.token}`,
        },
      }).toPromise();*/
      const TempItem: ItemDTO = {
        name: tem.name,
        guestOS: tem.guestOS,
        //guestOS: '',
        resourceID: tem.resourceID,
        dateAdded: tem.dateAdded,
        //vendor: res.vendor,
        vendor: '',
        //datacenter: res.organization,
        datacenter: '',
        memorySize: tem.memorySize,
      }

      List.push(TempItem);
    }
    return List;
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
  resourceID: number;
  guestOS: string;
  dateAdded: string;
  vendor: string;
  datacenter: string;
  memorySize: number;
}

interface ItemRes {
  id: string;
  vcdId: string;
  name: string;
  organization: string;
  vendor: string;
  // unit: GHz
  cpu: number;
  // unit: GB
  mem: number;
  // unit: GB
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

interface ItemTemplate {
  name: string;
  resourceID: number;
}

