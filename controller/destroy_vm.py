'''
Written by Zhe Shen, 19-11-3
Destroy VM in vSphere.
'''

from __future__ import print_function
import atexit
from pyVim import connect
from pyVmomi import vim
import argparse


def get_args():

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
    
    parser.add_argument('-j', '--uuid',
                        action='store',
                        help='BIOS UUID of the virtual machine')
    
    parser.add_argument('-d', '--dns',
                        action='store',
                        help='DNS name of the virtual machine')
    
    parser.add_argument('-i', '--ip',
                        help='IP address of the virtual machine')
    
    parser.add_argument('-n', '--name',
                        help='VM name of the virtual machine')

    args = parser.parse_args()

    return args


def get_obj(content, vimtype, name):

    obj = None
    container = content.viewManager.CreateContainerView(
        content.rootFolder, vimtype, True)
    for c in container.view:
        if name:
            if c.name == name:
                obj = c
                break
        else:
            obj = c
            break

    container.Destroy()
    return obj


args = get_args()
si = None
try:
    si = connect.SmartConnectNoSSL(host=args.host,
                                   user=args.user,
                                   pwd=args.password,
                                   port=args.port)
    atexit.register(connect.Disconnect, si)
except:
    print("Failed to connect")
    exit()

vm = None
if args.name:
    vm = get_obj(si.content, [vim.VirtualMachine], args.name)
elif args.uuid:
    vm = si.content.searchIndex.FindByUuid(None, args.uuid, True, False)
elif args.dns:
    vm = si.content.searchIndex.FindByDnsName(None, args.name, True)
elif args.ip:
    vm = si.content.searchIndex.FindByIp(None, args.ip, True)
else:
    print("Lack identifier of VM.")
    exit()

print("Found VM: {0}".format(vm.name))
print("The current power state is: {0}".format(vm.runtime.powerState))
if format(vm.runtime.powerState) == "poweredOn":
    print("Attempting to power off {0}".format(vm.name))
    task = vm.PowerOffVM_Task()
    print("{0}".format(task.info.state))

print("Destroying VM from vSphere.")
task = vm.Destroy_Task()
print("Done.")
