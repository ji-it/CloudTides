#!/usr/bin/python
'''
Written by Gaurav Dogra
Github: https://github.com/dograga

Script to extract cpu usage of esxhosts on vcenter for last 1 hour with multithreading
'''
import atexit
from pyVmomi import vim
from pyVim.connect import SmartConnect, Disconnect, SmartConnectNoSSL
import time
import datetime
from pyVmomi import vmodl
from threading import Thread
import ssl
from pytz import timezone
import argparse

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

    parser.add_argument('-i', '--info',
                        required=True,
                        action='store',
                        help='cpu, memmory or disk')

    parser.add_argument('--no-ssl',
                        action='store_true',
                        help='Skip SSL verification')
    
    args = parser.parse_args()

    if not args.password:
        args.password = getpass.getpass(
            prompt='Enter password')

    return args


class perfdata():
    def metricvalue(self,item,depth):
        maxdepth=10
        if hasattr(item, 'childEntity'):
            if depth > maxdepth:
                return 0
            else:
                item = item.childEntity
                item=self.metricvalue(item,depth+1)
        return item

    def run(self,content,vihost,item):
        output=[]
        try:
            perf_dict = {}
            perfManager = content.perfManager
            perfList = content.perfManager.perfCounter
            for counter in perfList: #build the vcenter counters for the objects
                counter_full = "{}.{}.{}".format(counter.groupInfo.key,counter.nameInfo.key,counter.rollupType)
                perf_dict[counter_full] = counter.key
            #print(perf_dict)
            cst_tz = timezone('Asia/Shanghai')
            for name in item:
                print("****************************************************************************")
                counter_name = name + '.usage.average'
                counterId = perf_dict[counter_name]
                metricId = vim.PerformanceManager.MetricId(counterId=counterId, instance="")
                timenow = datetime.datetime.now()
                #print(timenow)
                startTime = timenow - datetime.timedelta(hours=1)
                endTime = timenow
                search_index = content.searchIndex
                host = search_index.FindByDnsName(dnsName=vihost, vmSearch=False)
                query = vim.PerformanceManager.QuerySpec(entity=host,metricId=[metricId],intervalId=20,startTime=startTime,endTime=endTime)
                stats = perfManager.QueryPerf(querySpec=[query])
                count = 0
                for val in stats[0].value[0].value:
                    perfinfo={}
                    val=float(val/100)
                    perfinfo['timestamp']=startTime + datetime.timedelta(seconds=count*20)
                    perfinfo['hostname']=vihost
                    perfinfo['value']=val
                    output.append(perfinfo)
                    count+=1
                for out in output:
	                print("Hostname:{} TimeStamp: {} {} Usage: {}".format(out['hostname'],out['timestamp'],name,out['value']))
        except vmodl.MethodFault as e:
            print("Caught vmodl fault : " + e.msg)
            return 0
        except Exception as e:
            print("Caught exception : " + str(e))
            return 0

def main():
    
    try:
        args = get_args()
        if args.info != 'cpu' and args.info != 'memory' and args.info != 'disk':
            print("Wrong query info!")
            exit(0)
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
        perf=perfdata()
        for child in content.rootFolder.childEntity:
            datacenter=child
            hostfolder= datacenter.hostFolder
            hostlist=perf.metricvalue(hostfolder,0)
            for hosts in hostlist:
                esxhosts=hosts.host
                for esx in esxhosts:
                    summary=esx.summary
                    esxname=summary.config.name
                    p = Thread(target=perf.run, args=(content,esxname,[args.info]))
                    p.start()
    except:
        print("Failed to connect")
    

# start
if __name__ == "__main__":
    main()

