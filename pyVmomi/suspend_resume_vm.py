'''
Copyright 2013-2014 Reubenur Rahman
All Rights Reserved
@author: reuben.13@gmail.com
'''

import atexit
import argparse
import sys
import time
import ssl
import getpass

from pyVmomi import vim, vmodl
from pyVim import connect
from pyVim.connect import Disconnect, SmartConnect, GetSi, SmartConnectNoSSL


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
                        help='Name of the VM you wish to make')

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
    """
     Get the vsphere object associated with a given text name
    """    
    obj = None
    container = content.viewManager.CreateContainerView(content.rootFolder, vimtype, True)
    for c in container.view:
        if c.name == name:
            obj = c
            break
    return obj
'''
def wait_for_task(task, raiseOnError=True, si=None, pc=None):
    if si is None:
        si = GetSi()

    if pc is None:
        sc = si.RetrieveContent()
        pc = sc.propertyCollector

    # First create the object specification as the task object.
    objspec = vmodl.query.PropertyCollector.ObjectSpec()
    objspec.SetObj(task)

    # Next, create the property specification as the state.
    propspec = vmodl.query.PropertyCollector.PropertySpec()
    propspec.SetType(vim.Task);
    propspec.SetPathSet(["info.state"]);
    propspec.SetAll(True)

    # Create a filter spec with the specified object and property spec.
    filterspec = vmodl.query.PropertyCollector.FilterSpec()
    filterspec.SetObjectSet([objspec])
    filterspec.SetPropSet([propspec])

    # Create the filter
    filter = pc.CreateFilter(filterspec, True)
    
    # Loop looking for updates till the state moves to a completed state.
    taskName = task.GetInfo().GetName()
    update = pc.WaitForUpdates(None)
    state = task.GetInfo().GetState()
    while state != vim.TaskInfo.State.success and \
            state != vim.TaskInfo.State.error:
        if (state == 'running') and (taskName.info.name != "Destroy"):
            # check to see if VM needs to ask a question, thow exception
            vm = task.GetInfo().GetEntity()
            if vm is not None and isinstance(vm, vim.VirtualMachine):
                qst = vm.GetRuntime().GetQuestion()
            if qst is not None:
                raise Exception("Task blocked, User Intervention required")
      
    update = pc.WaitForUpdates(update.GetVersion())
    state = task.GetInfo().GetState()
         
    filter.Destroy()
    if state == "error" and raiseOnError:
        raise task.GetInfo().GetError()
      
    return state
'''

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
                invoke_and_track(vm. PowerOff)
        
        elif args.operation == 'suspend':
            if not force:
                invoke_and_track(vm.StandbyGuest)
            else:
                invoke_and_track(vm. Suspend)

        else:
            print("Wrong operation!")
            exit(0)
                
        #wait_for_task(task, si)        
        
    except:
        raise
    
# Start program
if __name__ == "__main__":
    main()