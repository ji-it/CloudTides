'''
Writted by Zhe Shen, 19-11-3
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


'''
def metricvalue(item, depth):
    maxdepth = 10
    if hasattr(item, 'childEntity'):
        if depth > maxdepth:
            return 0
        else:
            item = item.childEntity
            item = metricvalue(item, depth + 1)
    return item

def run(content, vihost, item, time, cur):
    output = []
        
    perf_dict = {}
    perfManager = content.perfManager
    perfList = content.perfManager.perfCounter
    for counter in perfList: #build the vcenter counters for the objects
        counter_full = counter.groupInfo.key + '.' + counter.nameInfo.key + '.' + counter.rollupType
        perf_dict[counter_full] = counter.key
    #print(perf_dict)
    #cst_tz = timezone('Asia/Shanghai')
    for name in item:
        print("****************************************************************************")
        info_name = None
        if name != 'power':
            info_name = name + '.usage.average'
        else:
            info_name = name + '.power.average'
        counterId = perf_dict[info_name]
        metricId = vim.PerformanceManager.MetricId(counterId=counterId)
        timenow = datetime.datetime.now()
        #print(timenow)
        startTime = timenow - datetime.timedelta(minutes=time)
        endTime = timenow
        search_index = content.searchIndex
        host = search_index.FindByDnsName(dnsName=vihost, vmSearch=False)
        query = vim.PerformanceManager.QuerySpec(entity=host,metricId=[metricId],intervalId=20,startTime=startTime,endTime=endTime)
        stats = perfManager.QueryPerf(querySpec=[query])
        count = 0
        for val in stats[0].value[0].value:
            perfinfo = {}
            val = float(val/100)
            perfinfo['timestamp'] = startTime + datetime.timedelta(seconds=count*20)
            perfinfo['hostname'] = vihost
            perfinfo['value'] = val
            output.append(perfinfo)
            count += 1
            cur.execute("INSERT INTO usage  (time, data) VALUES (%s, %s)", (perfinfo['timestamp'], val))
        for out in output:
	        print("Hostname:{} TimeStamp: {} {} Usage: {}".format(out['hostname'], out['timestamp'], name, out['value']))
'''


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
                # total_ram = host.hardware.memorySize / GBFACTOR
                #  total_cpu = round(((host.hardware.cpuInfo.hz / 1e+9) * host.hardware.cpuInfo.numCpuCores), 0)
                current_ram = float(host.summary.quickStats.overallMemoryUsage / 1024.0)
                current_cpu = float(host.summary.quickStats.overallCpuUsage / 1024.0)
                data['current_cpu'] = current_cpu
                data['current_ram'] = current_ram

                requests.post(FULL_HOSTNAME + "/api/usage/updatehost/", data=json.dumps(data))
                # print(data)

    # print(data)


# start
if __name__ == "__main__":
    main()
