import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from '@tide-environments/environment';
import { LoginService } from '../login/login.service';
import { TEMPLATEVM_PATH, TEMPLATEVM_PATH_NAME, TEMPLATE_PATH, VCD_URL_PATH} from '@tide-shared/config/path';

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
      const res = await this.http.get<ItemRes>(environment.apiPrefix + VCD_URL_PATH + `/` + tem.resourceID, {
        headers: {
          Authorization: `Bearer ${this.loginService.token}`,
        },
      }).toPromise();
      const TempItem: ItemDTO = {
        id: tem.id,
        name: tem.name,
        guestOS: tem.guestOS,
        //guestOS: '',
        resourceID: tem.resourceID,
        dateAdded: tem.dateAdded,
        vendor: res.vendor,
        //vendor: '',
        datacenter: res.datacenter,
        //datacenter: '',
        memorySize: tem.memorySize,
        provisionedSpace: tem.provisionedSpace,
        vcpu: tem.vcpu,
      }

      List.push(TempItem);
    }
    return List;
  }

  async getResList() {
    const res = await this.http.get<ItemRes[]>(environment.apiPrefix + VCD_URL_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const ResourceObject : Object = {};
    for (let item of res) {
      ResourceObject[item.datacenter] = item.id;
    }
    return ResourceObject;
  }

  async getVMList(id: number) {
    const TemplateList = await this.http.get<ItemVM[]>(environment.apiPrefix + TEMPLATE_PATH + `/` + TEMPLATEVM_PATH_NAME + `/` + id, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const List: ItemVM[] = [];
    for (const tem of TemplateList) {
      const TempItem: ItemVM = {
        id: tem.id,
        name: tem.name,
        vcpu: tem.vcpu,
        vmem: tem.vmem,
        disk: tem.disk,
        ports: tem.ports,
      }

      List.push(TempItem);
    }
    return List;
  }

  async getTemplateList() {
    const TemplateList = await this.http.get<ItemDTO[]>(environment.apiPrefix + TEMPLATE_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const TemplateObject : Object = {};
    for (let item of TemplateList){
      TemplateObject[item.name] = item.id;
    }
    return TemplateObject;
  }

  addItem(payload: ItemPayload) {
    const body = {
      ...payload,
    };
    return this.http.post<any>(environment.apiPrefix + TEMPLATE_PATH, body, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise().then(() => {
      return Promise.resolve();
    }, (errResp) => {
      return Promise.reject(`HTTP ${errResp.status}: ${errResp.error.message}`);
    });
  }

  addItemVM(payload: ItemPayloadVM) {
    const body = {
      ...payload,
    };
    return this.http.post<any>(environment.apiPrefix + TEMPLATEVM_PATH, body, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise().then(() => {
      return Promise.resolve();
    }, (errResp) => {
      return Promise.reject(`HTTP ${errResp.status}: ${errResp.error.message}`);
    });
  }

  editItem(id: string, payload: ItemPayload) {
    return this.http.put<ItemDTO>(`${this.prefix}/${id}`, payload).pipe(
      map(mapItem),
    );
  }

  async removeItem(id: string) {
    await this.http.delete<any>(environment.apiPrefix + TEMPLATE_PATH + `/` + id, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
    }, }).toPromise().then(
      () => {
        return Promise.resolve();
      }, (errResp) => {
        return Promise.reject(`${errResp.message}`);
      },
    );
  }

  async removeItemVM(id: number) {
    await this.http.delete<any>(environment.apiPrefix + TEMPLATEVM_PATH + `/` + id, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
    }, }).toPromise().then(
      () => {
        return Promise.resolve();
      }, (errResp) => {
        return Promise.reject(`${errResp.message}`);
      },
    );
  }
}

// Raw
interface ItemDTO {
  id: number;
  name: string;
  resourceID: number;
  guestOS: string;
  dateAdded: string;
  vendor: string;
  datacenter: string;
  memorySize: number;
  provisionedSpace: number;
  vcpu: number;
}

interface ItemRes {
  id: string;
  vcdId: string;
  datacenter: string;
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
  name: string,
  os: string,
  source: string,
  compat: string,
  space: number,
  memsize: number,
  vcpu: number,
  vmName: string,
  resourceID: number,
}

export interface ItemPayloadVM {
  name: string,
  disk: number,
  vmem: number,
  vcpu: number,
  templateID: number,
  ports: string,
}

export interface ItemVM {
  id: number,
  name: string,
  vcpu: number,
  vmem: number,
  disk: number,
  ports: string,
}

interface ItemTemplate {
  name: string;
  resourceID: number;
}

