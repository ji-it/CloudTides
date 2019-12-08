'''
Written by Zhe Shen, 19-11-2
Deploy a VM from a specified template.
'''

from pyVmomi import vim
from pyVim.connect import SmartConnect, SmartConnectNoSSL, Disconnect
import atexit
import argparse
import getpass
import random
import string
import time
import psycopg2
import os
import requests
import json
from tools import tasks # pylint: disable=import-error

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
    parser.add_argument('-v', '--vm-name',
                        required=True,
                        action='store',
                        help='Name of the VM you wish to make')
    '''
    parser.add_argument('--no-ssl',
                        action='store_true',
                        help='Skip SSL verification')

    parser.add_argument('--template',
                        required=True,
                        action='store',
                        help='Name of the template/VM \
                            you are cloning from')

    parser.add_argument('--datacenter-name',
                        required=False,
                        action='store',
                        default=None,
                        help='Name of the Datacenter you\
                            wish to use. If omitted, the first\
                            datacenter will be used.')

    parser.add_argument('--cluster-name',
                        required=False,
                        action='store',
                        default=None,
                        help='Name of the cluster you wish the VM to\
                            end up on. If left blank the first cluster found\
                            will be used')

    parser.add_argument('--resource-pool',
                        required=False,
                        action='store',
                        default=None,
                        help='Resource Pool to use. If left blank the first\
                            resource pool found will be used')

    parser.add_argument('--power-on',
                        dest='power_on',
                        action='store_true',
                        help='power on the VM after creation')
    
    
    args = parser.parse_args()

    if not args.password:
        args.password = getpass.getpass(
            prompt='Enter password')

    return args


def wait_for_task(task):
    ''' wait for a vCenter task to finish '''
    task_done = False
    while not task_done:
        if task.info.state == 'success':
            return task.info.result

        if task.info.state == 'error':
            print("there was an error")
            task_done = True


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


def clone_vm(content, template, si, datacenter_name, username, password,
        cluster_name, resource_pool, power_on, host_address, tem_name):

    # if none get the first one
    datacenter = get_obj(content, [vim.Datacenter], datacenter_name)
    destfolder = datacenter.vmFolder

    # if None, get the first one
    cluster = get_obj(content, [vim.ClusterComputeResource], cluster_name)

    if resource_pool:
        resource_pool = get_obj(content, [vim.ResourcePool], resource_pool)
    else:
        resource_pool = cluster.resourcePool

    vmconf = vim.vm.ConfigSpec()
    now = int(round(time.time()*1000))
    now02 = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(now/1000))
    my_annotation = "Here is a BOINC-deployed VM created by ProjectTides at " + now02 + " with BOINC unstarted"
    vmconf.annotation = my_annotation

    # set relospec
    relospec = vim.vm.RelocateSpec()
    relospec.datastore = None
    relospec.pool = resource_pool

    clonespec = vim.vm.CloneSpec(powerOn=power_on, template=False, location=relospec)
    clonespec.config = vmconf
    vm_name = ''.join(random.sample(string.ascii_letters + string.digits, 8))
    print("cloning VM...")
    task = template.Clone(folder=destfolder, name=vm_name, spec=clonespec)
    wait_for_task(task)
    time.sleep(60)
    VM = get_obj(content, [vim.VirtualMachine], vm_name)
    #print(VM)
    print(VM.summary.guest.ipAddress)
    send_account(host_address, VM.summary.guest.ipAddress, tem_name, username, password)

    spec = vim.vm.ConfigSpec()
    old_ann = VM.summary.config.annotation
    now = int(round(time.time()*1000))
    now02 = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(now/1000))
    annotation_add = "started at " + now02
    print(old_ann)
    spec.annotation = old_ann.replace("unstarted", annotation_add)

    task = VM.ReconfigVM_Task(spec)
    tasks.wait_for_tasks(si, [task])
    print("Done.")


def send_account(host_address, ip_address, template, username, password):
    conn=psycopg2.connect(database="Tides",user="postgres",
            password="t6bB2T5KoQuPq6DrpWxJa3rYKVjIpOCtVSrKyBMB8PHcMShkidcQo8Kjn1lcXswB",host="10.11.16.83",port="30123")
    cur=conn.cursor()

    cur.execute("SELECT policy_id FROM resource_resource WHERE host_address = %s", (host_address,))
    pid = cur.fetchone()

    cur.execute("SELECT account_type, project_id, username, password FROM policy_policy WHERE id = %s", str(pid[0]))
    result = cur.fetchone()

    cur.execute("SELECT url FROM projects_projects WHERE id = %s", str(result[1]))
    url = cur.fetchone()

    cur.execute("SELECT password FROM template_template WHERE name = %s", (template,))
    pwd = cur.fetchone()
    filename = 'account.txt'
    with open(filename, 'w') as f:
        f.write(str(result[0]) + '\n')
        f.write(url[0] + '\n')
        f.write(result[2] + '\n')
        f.write(result[3])
    os.system('sshpass -p ' + pwd[0] + ' scp ' + filename + ' root@' + ip_address + ':/var/lib/boinc')
    os.system('sshpass -p ' + pwd[0] + ' scp run_boinc root@' + ip_address + ':/var/lib/boinc')
    os.system('python execute_program.py -s ' + host_address + ' -u ' + username + ' -p ' + password + ' -S -i ' + ip_address +
                ' -r root -w ' + pwd[0] + ' -l /bin/chmod -f "777 /var/lib/boinc/run_boinc"')
    os.system('python execute_program.py -s ' + host_address + ' -u ' + username + ' -p ' + password + ' -S -i ' + ip_address +
                ' -r root -w ' + pwd[0] + ' -l /var/lib/boinc/run_boinc -f None')
    '''
    data = {}
    data['host'] = host_address
    data['username'] = username
    data['password'] = password
    data['ip_address'] = ip_address
    requests.post("http://192.168.56.1:8000/api/usage/setannotation/", data=json.dumps(data))
    '''



def main():
    
    args = get_args()

    # connect this thing
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

    content = si.RetrieveContent()
    template = None

    template = get_obj(content, [vim.VirtualMachine], args.template)
    
    if template:
        clone_vm(content, template, si, args.datacenter_name, args.user, args.password,
            args.cluster_name, args.resource_pool, args.power_on, args.host, args.template)
    else:
        print("template not found")


# start this thing
if __name__ == "__main__":
    main()
