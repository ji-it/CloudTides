'''
Written by Zhe Shen, 19-11-3
Suspend and resume VM.
'''

import atexit
import argparse
import ssl
import getpass

from pyVmomi import vim
from pyVim import connect
from pyVim.connect import Disconnect, SmartConnect, SmartConnectNoSSL


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

    parser.add_argument('-v', '--vm-name',
                        required=True,
                        action='store',
                        help='Name of the VM you wish to operate on')

    parser.add_argument('--no-ssl',
                        action='store_true',
                        help='Skip SSL verification')

    parser.add_argument('--operation',
                        required=True,
                        action='store',
                        help='start, suspend, or stop')

    parser.add_argument('-f', '--force',
                        required=False,
                        action='store',
                        default=None)
    
    args = parser.parse_args()

    if not args.password:
        args.password = getpass.getpass(
            prompt='Enter password')

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


def invoke_and_track(func, *args, **kw):
    try :
        task = func(*args, **kw)
        #wait_for_task(task)
    except:
        raise


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
        print("Connection failed.")
    # disconnect this thing
    atexit.register(Disconnect, si)

    print("Connected to VCENTER SERVER !")
        
    content = si.RetrieveContent()

    if args.operation == 'stop' or args.operation == 'suspend':
        force = args.force
    
    vm = get_obj(content, [vim.VirtualMachine], args.vm_name)

    #current_state = vm.runtime.powerState
        
    if args.operation == 'start':
        invoke_and_track(vm.PowerOn, None)

    elif args.operation == 'stop':
        if not force:
            invoke_and_track(vm.ShutdownGuest)
        else:
            invoke_and_track(vm.PowerOff)
        
    elif args.operation == 'suspend':
        if not force:
            invoke_and_track(vm.StandbyGuest)
        else:
            invoke_and_track(vm.Suspend)

    else:
        print("Wrong operation!")
        exit()
                
    #wait_for_task(task, si)        
        
    
# Start program
if __name__ == "__main__":
    main()