'''
Query CPU, memory, disk and power usage of ESXi hosts.
'''

import ssl
import argparse
import atexit
from pyVim.connect import SmartConnect, Disconnect, SmartConnectNoSSL
import getpass
import requests
import json
from pyVmomi import vim
from config import FULL_HOSTNAME

GBFACTOR = float(1 << 30)
requests.adapters.DEFAULT_RETRIES = 5



def get_args():
    """ Get arguments from CLI """
    parser = argparse.ArgumentParser(
        description='Arguments for talking to vCenter')

    parser.add_argument('-s', '--host',
                        required=True,
                        action='store',
                        help='vSpehre service to connect to')

    parser.add_argument('-o', '--port',
                        type=int,
                        default=443,
                        action='store',
                        help='Port to connect on')

    parser.add_argument('-u', '--user',
                        required=True,
                        action='store',
                        help='Username to use')

    parser.add_argument('-p', '--password',
                        required=False,
                        action='store',
                        help='Password to use')
    '''
    parser.add_argument('-i', '--info',
                        required=True,
                        action='store',
                        help='cpu, mem or disk')
    '''
    parser.add_argument('-n', '--name',
                        required=False,
                        action='store',
                        help='resource name to query')
    
    parser.add_argument('-r', '--resource-pool',
                        action='store_true',
                        help='whether is resource pool')

    parser.add_argument('--no-ssl',
                        action='store_true',
                        help='Skip SSL verification')

    args = parser.parse_args()

    if not args.password:
        args.password = getpass.getpass(
            prompt='Enter password')

    return args


def get_all_objs(content, vimtype, folder=None, recurse=True):
    if not folder:
        folder = content.rootFolder

    obj = {}
    container = content.viewManager.CreateContainerView(folder, vimtype, recurse)
    for managed_object_ref in container.view:
        obj.update({managed_object_ref: managed_object_ref.name})
    return obj


def main():
    try:
        args = get_args()
        si = None
        if args.no_ssl:
            si = SmartConnectNoSSL(
                host=args.host,
                user=args.user,
                pwd=args.password,
                port=args.port)
        else:
            si = SmartConnect(
                host=args.host,
                user=args.user,
                pwd=args.password,
                port=args.port)
    except:
        print("Failed to connect")
        exit()
    # disconnect this thing
    atexit.register(Disconnect, si)
    content = si.RetrieveContent()

    objs = None
    data = {}
    data['HostAddress'] = args.host
    data['Name'] = args.name
    if args.resource_pool:
        objs = get_all_objs(content, [vim.ResourcePool])
        for pool in objs:
            if pool.name == args.name:
                info = pool.runtime
                data['CurrentCPU'] = info.cpu.overallUsage
                data['TotalCPU'] = info.cpu.maxUsage
                data['CurrentRAM'] = info.memory.overallUsage / (1024.0*1024.0)
                data['TotalRAM'] = info.memory.maxUsage / (1024.0*1024.0)
                break
    else:
        objs = get_all_objs(content, [vim.ClusterComputeResource])
        for cluster in objs:
            if cluster.name == args.name:
                summary = cluster.GetResourceUsage()
                data['CurrentCPU'] = summary.cpuUsedMHz / 1000.0
                data['TotalCPU'] = summary.cpuCapacityMHz / 1000.0
                data['CurrentRAM'] = summary.memUsedMB / 1024.0
                data['TotalRAM'] = summary.memCapacityMB / 1024.0
                break
    
    print(data)
    headers = {'Content-type': 'application/json'}
    res = requests.put(FULL_HOSTNAME + "/v1/usage/update_resource/", data=json.dumps(data), headers=headers)
    print(res)


# start
if __name__ == "__main__":
    main()
