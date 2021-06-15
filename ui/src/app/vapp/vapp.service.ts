import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { environment } from '@tide-environments/environment';
import { PORT_PATH_NAME, TEMPLATE_PATH, VAPP_PATH, VAPP_PATH_NAME, VCD_URL_PATH, VENDOR_PATH, VM_PATH_NAME} from '@tide-config/path';
import { LoginService } from '../login/login.service';
import toFixed from 'accounting-js/lib/toFixed.js';

@Injectable()
export class VappService {

  constructor(
    private readonly http: HttpClient,
    private readonly loginService: LoginService,
  ) {
  }

  private prefix = `${environment.apiPrefix}/computeResource`;

  async getList() {
    const list = await this.http.get<Item[]>(environment.apiPrefix + VAPP_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const vapps: Item[] = [];
    for (const vapp of list) {
      const vappItem: Item = {
        id: vapp.id,
        name: vapp.name,
        vendor: vapp.vendor,
        template: vapp.template,
        datacenter: vapp.datacenter,
        ipaddress: vapp.ipaddress,
        status: vapp.status,
        poweredOn: vapp.poweredOn,
      };
      vapps.push(vappItem);
    }
    return vapps;
  }

  async getVMList(id: number) {
    const VMList = await this.http.get<ItemVM[]>(environment.apiPrefix + `/` + VAPP_PATH_NAME + `/` + VM_PATH_NAME + `/` + id, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const List: ItemVM[] = [];
    for (const tem of VMList) {
      const TempItem: ItemVM = {
        id: tem.id,
        name: tem.name,
        vcpu: tem.vcpu,
        vmem: tem.vmem,
        disk: tem.disk,
        IPAddress: tem.IPAddress,
        externalIPAddress: tem.externalIPAddress,
        username: tem.username,
        password: tem.password,
        usedMoney: tem.usedMoney,
        status: tem.status,
      }

      List.push(TempItem);
    }
    return List;
  }

  async getVendorList() {
    const VendorList = await this.http.get<ItemVendor[]>(environment.apiPrefix + VENDOR_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const VendorObject : Object = {};
    for (let item of VendorList){
      VendorObject[item.name] = item.url;
    }
    return VendorObject;
  }

  async getTemplateList() {
    const TemplateList = await this.http.get<ItemTemplate[]>(environment.apiPrefix + TEMPLATE_PATH, {
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

  async getResourceList() {
    const list = await this.http.get<ItemResource[]>(environment.apiPrefix + VCD_URL_PATH, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    const usage: ItemResource[] = [];
    for (const resource of list) {
      const rawUsage = await this.http.get<ItemUsage>(`${environment.apiPrefix}/usage/${resource.id}`, {
        headers: {
          Authorization: `Bearer ${this.loginService.token}`,
        },
      }).toPromise();
      const resourceItem: ItemResource = {
        id: resource.id,
        vcdId: resource.vcdId,
        name: rawUsage.name,
        organization: resource.organization,
        vendor: resource.vendor,
        cpu: rawUsage.totalCPU / 1000,
        mem: rawUsage.totalRAM / 1024,
        disk: rawUsage.totalDisk / 1024,
        resType: resource.resType,
        usage: {
          'cpu%': toFixed(rawUsage.percentCPU * 100, 2),
          'mem%': toFixed(rawUsage.percentRAM * 100, 2),
          'disk%': toFixed(rawUsage.percentDisk * 100, 2),
        },
      };
      usage.push(resourceItem);
    }
    const ResourceObject : Object = {};
    for (let item of usage) {
      ResourceObject[item.name] = item.name;
    }
    return ResourceObject;
  }

  async getPortList(id: number) {
    const PortList = await this.http.get<ItemPort[]>(environment.apiPrefix + `/` + VM_PATH_NAME + `/` + PORT_PATH_NAME + `/` + id, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise();
    return PortList;
  }

  addItem(payload: ItemPayload) {
    const body = {
      ...payload,
    };
    return this.http.post<any>(environment.apiPrefix + VAPP_PATH, body, {
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
    await this.http.delete<any>(environment.apiPrefix + VAPP_PATH + `/` + id, {
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

  async contributeResource(id: string): Promise<ContributeResp> {
    let response = null;
    await this.http.put(environment.apiPrefix + `/resource/contribute/${id}`, null, {
      headers: {
      Authorization: `Bearer ${this.loginService.token}`,
    }, }).toPromise().then(
      (resp) => {
        response = resp;
        return Promise.resolve();
      }, (errResp) => {
        return Promise.reject(`${errResp.message}`);
      },
    );
    return response;
  }

  async activateResource(id: string): Promise<ActivateResp> {
    let response = null;
    await this.http.put<ActivateResp>(environment.apiPrefix + `/resource/activate/${id}`, null, {
      headers: {
        Authorization: `Bearer ${this.loginService.token}`,
      },
    }).toPromise().then(
      (resp) => {
        response = resp;
        return Promise.resolve();
      }, (errResp) => {
        return Promise.reject(`${errResp.message}`);
      },
    );
    return response;
  }
}

interface ItemDTO {
  id: number;
  name: string;
  vendor: string;
  datacenter: string;
  template: string;
  ipaddress: string;
  status: string;
  poweredOn: boolean;
}

interface ContributeResp {
  message: string;
  contributed: boolean;
}

interface ActivateResp {
  message: string;
  activated: boolean;
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
  template: number;
  vendor: string;
  datacenter: string;
}

interface ItemVendor {
  id: number;
  name: string;
  url: string;
  vendorType: string;
  version: string;
}

interface ItemTemplate {
  id: number;
  name: string;
  resourceID: number;
}

export interface ItemResource {
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
  resType: string;
  usage: {
    'cpu%': number;
    'mem%': number;
    'disk%': number;
  }
}

interface ItemUsage {
  currentCPU: number;
  totalCPU: number;
  currentRAM: number;
  totalRAM: number;
  currentDisk: number;
  totalDisk: number;
  percentCPU: number;
  percentRAM: number;
  percentDisk: number;
  name: string;
}

export interface ItemVM {
  id: number,
  name: string,
  vcpu: number,
  vmem: number,
  disk: number,
  IPAddress: string,
  externalIPAddress: string,
  username: string,
  password: string,
  usedMoney: number,
  status: string,
}

export interface ItemPort {
  id: number,
  port: number,
  url: string,
}


export type ItemV = ItemVendor;

export type Item = ItemDTO;
