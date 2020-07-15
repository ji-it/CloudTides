'''
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
from tools import tasks  # pylint: disable=import-error
from config import BASE_DIR, DATABASES
import paramiko
import re


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
                        required=True,
                        action='store',
                        help='Name of the Host on which you wish to deploy VM')
    
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


def clone_vm(content, template, si, datacenter_name, username, password, name,
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
    now = int(round(time.time() * 1000))
    now02 = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(now / 1000))
    my_annotation = "Here is a BOINC-deployed VM created by ProjectTides at " + now02 + " with BOINC unstarted"
    vmconf.annotation = my_annotation

    # set relospec
    relospec = vim.vm.RelocateSpec()
    relospec.datastore = None
    relospec.pool = resource_pool

    clonespec = vim.vm.CloneSpec(powerOn=power_on, template=False, location=relospec)
    clonespec.config = vmconf
    vm_name = 'tides_worker-' + ''.join(random.sample(string.ascii_letters + string.digits, 8))
    print("cloning VM...")
    task = template.Clone(folder=destfolder, name=vm_name, spec=clonespec)
    
    # time.sleep(120)
    VM = wait_for_task(task)
    print(VM)
    while VM.summary.guest.ipAddress is None:
        pass
    print(VM.summary.guest.ipAddress)
    send_account(host_address, name, VM.summary.guest.ipAddress, tem_name, username, password)

    spec = vim.vm.ConfigSpec()
    old_ann = VM.summary.config.annotation
    now = int(round(time.time() * 1000))
    now02 = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(now / 1000))
    annotation_add = "started at " + now02
    print(old_ann)
    spec.annotation = old_ann.replace("unstarted", annotation_add)

    task = VM.ReconfigVM_Task(spec)
    tasks.wait_for_tasks(si, [task])
    print("Done.")


def send_account(host_address, name, ip_address, template, username, password):
    db = DATABASES['default']['NAME']
    user = DATABASES['default']['USER']
    dbpwd = DATABASES['default']['PASSWORD']
    host = DATABASES['default']['HOST']
    port = DATABASES['default']['PORT']
    conn = psycopg2.connect(database=db, user=user, password=dbpwd, host=host, port=port)
    cur = conn.cursor()

    cur.execute("SELECT policy_ref FROM resources WHERE host_address = %s AND name = %s", (host_address, name,))
    pid = cur.fetchone()

    cur.execute("SELECT account_type, project_ref, boinc_username, boinc_password FROM policies WHERE id = %s", str(pid[0]))
    result = cur.fetchone()

    cur.execute("SELECT url FROM projects WHERE id = %s", str(result[1]))
    url = cur.fetchone()

    # cur.execute("SELECT password FROM template_template WHERE name = %s", (template,))
    # pwd = cur.fetchone()
    pwd = "vmware"  # TODO: change to passwordless
    # path = os.path.join(BASE_DIR, 'controller')
    filename = 'account.txt'
    with open(filename, 'w') as f:
        f.write(str(result[0]) + '\n')
        f.write(url[0] + '\n')
        f.write(result[2] + '\n')
        f.write(result[3])
    # time.sleep(30)
    # run_boinc = 'run_boinc'
    path = os.path.join(BASE_DIR, 'controller')
    print("SSH Connecting...")
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    ssh.connect(ip_address, username = 'root', password = 'vmware', timeout = 5)

    while os.system("sshpass -p '" + pwd + "' scp " + filename + ' root@' + ip_address + ':/var/lib/boinc') != 0:
        time.sleep(5)
        continue
    while os.system("sshpass -p '" + pwd + "' scp " + path + '/run_boinc' + ' root@' + ip_address + ':/var/lib/boinc') != 0:
        time.sleep(5)
        continue
    '''
    os.system(
        'python3 ' + path + '/execute_program.py -s ' + host_address + ' -u ' + username + ' -p ' + password + ' -S -i ' + ip_address +
        " -r root -w '" + pwd + "' -l /bin/chmod -f '777 /var/lib/boinc/run_boinc'")
    os.system(
        'python3 ' + path + '/execute_program.py -s ' + host_address + ' -u ' + username + ' -p ' + password + ' -S -i ' + ip_address +
        " -r root -w '" + pwd + "' -l /var/lib/boinc/run_boinc -f None")
    '''
    
    cmd = 'cd /var/lib/boinc && chmod +x run_boinc && ./run_boinc'
    stdin, stdout, stderr = ssh.exec_command(cmd)
    ssh.close()
    '''
    master_ip = '10.185.143.234'
    JoinCommand = get_join_command(master_ip, ip_address)

    ssh.connect(ip_address, username = 'root', password = 'vmware', timeout = 5)
    cmd = 'kubeadm reset -f 2>&1 && %s 2>&1 \n' % JoinCommand

    stdin, stdout, stderr = ssh.exec_command(cmd)
    print(stdout.read())
    ssh.close()
    # print(stdin + '\n' + stdout + '\n' + stderr)
    '''
    # cur.execute("UPDATE resources SET status = 'contributing' WHERE host_address = %s AND name = %s", (host_address, name,))
    # conn.commit()

'''
def get_join_command(master_ip, node_ip):
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    ssh.connect(master_ip, username = 'root', password = 'vmware', timeout = 5)

    cmd = 'kubeadm token create --print-join-command 2>&1 \n'
    stdin, stdout, stderr = ssh.exec_command(cmd)
    JoinCommand = stdout.read().rstrip().decode('utf-8')

    ssh.close()
    if not re.findall("\-\-discovery\-token\-ca\-cert\-hash", JoinCommand):
        print("NOT a real join command")
        exit(-1)
    
    return JoinCommand
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
        clone_vm(content, template, si, args.datacenter_name, args.user, args.password, args.name,
                 args.cluster_name, args.resource_pool, args.power_on, args.host, args.template)
    else:
        print("template not found")


# start this thing
if __name__ == "__main__":
    main()
