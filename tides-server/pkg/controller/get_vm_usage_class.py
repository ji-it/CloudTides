from pyVim.connect import SmartConnectNoSSL, Disconnect
import atexit
import json
import re
import requests
import argparse
import getpass
from pyVmomi import vim
from config import FULL_HOSTNAME


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
    
    parser.add_argument('-n', '--name',
                        required=False,
                        action='store',
                        help='resource name to query')

    parser.add_argument('--no-ssl',
                        action='store_true',
                        help='Skip SSL verification')

    args = parser.parse_args()

    if not args.password:
        args.password = getpass.getpass(
            prompt='Enter password')

    return args


def printvminfo(vm_collect, vm, depth=1):
    """
    Print information for a particular virtual machine or recurse into a folder
    with depth protection
    """

    # if this is a group it will have children. if it does, recurse into them
    # and then return
    if hasattr(vm, 'childEntity'):
        if depth > 4:
            return
        vmlist = vm.childEntity
        for child in vmlist:
            printvminfo(vm_collect, child, depth + 1)
        return

    summary = vm.summary
    my_ann = summary.config.annotation
    create_time = ""
    boinc_time = "unstarted"
    if my_ann.find("Here is a BOINC-deployed VM created by ProjectTides") != -1:
        # print(my_ann)
        if my_ann.find("unstarted") != -1:
            searchObj = re.search('Here is a BOINC-deployed VM created by ProjectTides at (.*) with BOINC unstarted',
                                  my_ann)
            create_time = searchObj.group(1)
        else:
            searchObj = re.search(
                'Here is a BOINC-deployed VM created by ProjectTides at (.*) with BOINC started at (.*)', my_ann)
            create_time = searchObj.group(1)
            boinc_time = searchObj.group(2)

        vm_collect[summary.guest.ipAddress] = {
            "Name": summary.config.name,
            "CurrentRAM": float(summary.quickStats.guestMemoryUsage / 1024.0),
            "CurrentCPU": float(summary.quickStats.overallCpuUsage / 1000.0),
            # "vm_created_time": create_time,
            "Datacenter": summary.vm.parent.parent.name,
            "GuestOS": summary.config.guestFullName,
            "NumCPU": summary.config.numCpu,
            "HostName": summary.guest.hostName,
            "TotalRAM": float(summary.runtime.maxMemoryUsage / 1024.0),
            "TotalCPU": float(summary.runtime.maxCpuUsage / 1024.0),
            "PoweredOn": summary.runtime.powerState == "poweredOn",
            "BoincStartTime": boinc_time,
            # "ip_address": summary.guest.ipAddress
        }


def main():
    args = get_args()
    try:
        si = None
        si = SmartConnectNoSSL(
            host=args.host,
            user=args.user,
            pwd=args.password,
            port=args.port)
        atexit.register(Disconnect, si)
    except:
        print("Connection failed")
        exit()

    if si is None:
        exit()

    content = si.RetrieveContent()
    count = 0
    vm_collect = {}
    vm_collect["VMs"] = {}
    # vm_collect[args.host] = {}
    for child in content.rootFolder.childEntity:
        if hasattr(child, 'vmFolder'):
            datacenter = child
            vmfolder = datacenter.vmFolder
            vmlist = vmfolder.childEntity
            for vm in vmlist:
                if type(vm) == vim.VirtualMachine:
                    count = count + 1
                    if format(vm.runtime.powerState) == "poweredOn":
                        printvminfo(vm_collect['VMs'], vm)

    # print(vm_collect)
    vm_collect["total_vms"] = count
    vm_collect["Name"] = args.name
    print(vm_collect)
    headers = {'Content-type': 'application/json'}
    requests.post(FULL_HOSTNAME + "/v1/usage/addVM/", data=json.dumps(vm_collect), headers=headers)


if __name__ == "__main__":
    main()
