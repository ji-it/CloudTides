from pyVim.connect import SmartConnectNoSSL, Disconnect
import atexit

import json
import pyVmomi
from pyVmomi import vim, vmodl
import datetime

import re
import time
import requests
import argparse
import getpass

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
            printvminfo(vm_collect, child, depth+1)
        return
    
    summary = vm.summary
    my_ann = summary.config.annotation
    create_time = ""
    boinc_time = "unstarted"
    if my_ann.find("Here is a BOINC-deployed VM created by ProjectTides") != -1:
        print(my_ann)
        if my_ann.find("unstarted") != -1:
            searchObj = re.search('Here is a BOINC-deployed VM created by ProjectTides at (.*) with BOINC unstarted', my_ann)
            create_time = searchObj.group(1)
        else:
            searchObj = re.search('Here is a BOINC-deployed VM created by ProjectTides at (.*) with BOINC started at (.*)', my_ann)
            create_time = searchObj.group(1)
            boinc_time = searchObj.group(2)

        vm_collect[summary.guest.ipAddress] = {
            "Name": summary.config.name, 
            "Memory": summary.quickStats.guestMemoryUsage, 
            "CPU": summary.quickStats.overallCpuUsage, 
            "CreateTime": create_time,
            "BOINCTime": boinc_time
        }


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
#print(content)
vm_collect = {}
vm_collect[args.host] = {}
for child in content.rootFolder.childEntity:
    if hasattr(child, 'vmFolder'):
        datacenter = child
        vmfolder = datacenter.vmFolder
        vmlist = vmfolder.childEntity
        for vm in vmlist:
            if type(vm) != vim.Folder and format(vm.runtime.powerState) == "poweredOn":
                printvminfo(vm_collect, vm)
print(vm_collect)
requests.post("http://192.168.56.1:8000/api/usage/addvm/", data=json.dumps(vm_collect))