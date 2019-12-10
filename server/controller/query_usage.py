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

GBFACTOR = float(1 << 30)
requests.adapters.DEFAULT_RETRIES = 5

FULL_HOSTNAME = "http://localhost:8000"


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
                        help='Host name to query')

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

    children = content.rootFolder.childEntity
    for child in children:  # Iterate though DataCenters
        dc = child
        clusters = dc.hostFolder.childEntity
        for cluster in clusters:  # Iterate through the clusters in the DC
            hosts = cluster.host  # Variable to make pep8 compliance
            for host in hosts:  # Iterate through Hosts in the Cluster
                hostname = host.summary.config.name
                if hostname != args.name:
                    continue

                data = {}
                data['host_address'] = args.host
                data['host_name'] = hostname
                current_ram = float(host.summary.quickStats.overallMemoryUsage / 1024.0)
                current_cpu = float(host.summary.quickStats.overallCpuUsage / 1024.0)
                data['current_cpu'] = current_cpu
                data['current_ram'] = current_ram

                requests.post(FULL_HOSTNAME + "/api/usage/updatehost/", data=json.dumps(data))


# start
if __name__ == "__main__":
    main()
