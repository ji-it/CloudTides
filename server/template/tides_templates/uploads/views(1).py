# -*- coding: utf-8 -*-
from __future__ import unicode_literals, print_function

from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework import status
from .serializers import *

from .models import Resource
from pyVim.connect import SmartConnectNoSSL, Disconnect
import atexit
from rest_framework.response import Response
from rest_framework import status
import json
import pyVmomi
from pyVmomi import vim, vmodl
import datetime
from .utils import *

from tools import cli
from tools import tasks
import re
import time

GBFACTOR = float(1 << 30)


class ValidateResource(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']

        try:
            si = None
            si = SmartConnectNoSSL(
                host=host,
                user=username,
                pwd=password,
                port=443)
            atexit.register(Disconnect, si)
        except:
            return Response({'message': 'Failure'}, status=401)
        
        if si is None:
            return Response({'message': 'Failure'}, status=401)
        
        content = si.RetrieveContent()
        dic = []
        for datacenter in content.rootFolder.childEntity:
            dic.append(datacenter.name)
        return Response({'message': 'Success', 'datacenter names': dic}, status=status.HTTP_200_OK)


class AddResource(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        platform_type = data['platform_type']
        datacenter_name = data['datacenter']

        try:
            si = None
            si = SmartConnectNoSSL(
                host=host,
                user=username,
                pwd=password,
                port=443)
            atexit.register(Disconnect, si)
        except:
            return Response({'message': 'connection failure'}, status=401)
        
        if si is None:
            return Response({'message': 'connection failure'}, status=401)
        
        content = si.RetrieveContent()
        datacenters = get_all_objs(content, [vim.Datacenter])
        datacenter = None
        for dc in datacenters:
            if dc.name == datacenter_name:
                datacenter = dc
                break
        
        if datacenter is None:
            return Response({'message': 'datacenter not found'}, status=401)
        
        date_added = datetime.datetime.now()
        host_address = None
        total_ram = None
        total_cpu = None
        current_ram = None
        current_cpu = None
        is_active = True
        if hasattr(datacenter.hostFolder, 'childEntity'):
            hostFolder = datacenter.hostFolder
            computeResourceList = hostFolder.childEntity
            for computeResource in computeResourceList:
                hostList = computeResource.host
                for ho in hostList:
                    host_address = ho.name
                    total_ram = ho.hardware.memorySize/GBFACTOR
                    total_cpu = round(((ho.hardware.cpuInfo.hz/1e+9)*ho.hardware.cpuInfo.numCpuCores),0)
                    current_ram = float(ho.summary.quickStats.overallMemoryUsage/1024)
                    current_cpu = float(ho.summary.quickStats.overallCpuUsage/1024)
                    #print(current_cpu, current_ram, total_ram, total_cpu)
                    profile = Resource(username=username, password=password, date_added=date_added, host_address=host_address,
                                platform_type=platform_type, total_ram=total_ram, total_cpu=total_cpu, current_ram=current_ram,
                                current_cpu=current_cpu, is_active=is_active)
                    try:
                        profile.save()
                    except:
                        return Response({'message': 'resource already registered'}, status=200)

        return Response({'message': 'success'}, status=200)


class DeleteResource(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        platform_type = data['platform_type']
        datacenter_name = data['datacenter']

        try:
            si = None
            si = SmartConnectNoSSL(
                host=host,
                user=username,
                pwd=password,
                port=443)
            atexit.register(Disconnect, si)
        except:
            return Response({'message': 'connection failure'}, status=401)
        
        if si is None:
            return Response({'message': 'connection failure'}, status=401)
        
        content = si.RetrieveContent()
        datacenters = get_all_objs(content, [vim.Datacenter])
        datacenter = None
        for dc in datacenters:
            if dc.name == datacenter_name:
                datacenter = dc
                break
        
        if datacenter is None:
            return Response({'message': 'datacenter not found'}, status=401)

        if hasattr(datacenter.hostFolder, 'childEntity'):
            hostFolder = datacenter.hostFolder
            computeResourceList = hostFolder.childEntity
            for computeResource in computeResourceList:
                hostList = computeResource.host
                for ho in hostList:
                    host_address = ho.name
                    try:
                        Resource.objects.get(host_address=host_address).delete()
                    except:
                        return Response({'message': 'no matching objects'}, status=401)
        
        return Response({'message': 'success'}, status=200)

class DeployVMFromTemplate(APIView):
    def wait_for_task(self, task):
        ''' wait for a vCenter task to finish '''
        task_done = False
        while not task_done:
            if task.info.state == 'success':
                return task.info.result

            if task.info.state == 'error':
                print("there was an error")
                task_done = True

    def get_obj(self, content, vimtype, name):
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

    def clone_vm(self, content, template, vm_name, si, datacenter_name,
        cluster_name=None, resource_pool=None):

        # if none get the first one
        datacenter = self.get_obj(content, [vim.Datacenter], datacenter_name)
        destfolder = datacenter.vmFolder

        # if None, get the first one
        cluster = self.get_obj(content, [vim.ClusterComputeResource], cluster_name)

        if resource_pool:
            resource_pool = self.get_obj(content, [vim.ResourcePool], resource_pool)
        else:
            resource_pool = cluster.resourcePool

        vmconf = vim.vm.ConfigSpec()
        now = int(round(time.time()*1000))
        now02 = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(now/1000))
        my_annotation = "Here is a BOINC-deployed VM created by ProjectTides at " + now02 + " with BOINC unstarted"
        vmconf.annotation = my_annotation

        # set relospec
        relospec = vim.vm.RelocateSpec()
        # print(relospec)
        relospec.datastore = None
        relospec.pool = resource_pool

        clonespec = vim.vm.CloneSpec(powerOn=True, template=False, location=relospec)
        clonespec.config = vmconf

        print("cloning VM...")
        task = template.Clone(folder=destfolder, name=vm_name, spec=clonespec)
        self.wait_for_task(task)
        print("Done.")


    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        vm_name = data['vm_name']
        template = data['template']
        datacenter_name = data['datacenter_name']

        try:
            si = None
            si = SmartConnectNoSSL(
                host=host,
                user=username,
                pwd=password,
                port=443)
            atexit.register(Disconnect, si)
        except:
            return Response({'message': 'Failure'}, status=401)

        if si is None:
            return Response({'message': 'Failure'}, status=401)

        content = si.RetrieveContent()
        template = None
        template = self.get_obj(content, [vim.VirtualMachine], template)

        if template:
            self.clone_vm(content, template, vm_name, si, datacenter_name)
            return Response({'message': 'success'}, status=200)
        else:
            print("template not found")
            return Response({'message': 'no matching objects'}, status=401)

class DestroyVM(APIView):
    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        ip_address = data['ip_address']

        try:
            si = None
            si = SmartConnectNoSSL(
                host=host,
                user=username,
                pwd=password,
                port=443)
            atexit.register(Disconnect, si)
        except:
            return Response({'message': 'connection failure'}, status=401)
        
        if si is None:
            return Response({'message': 'connection failure'}, status=401)
        
        vm = None
        if ip_address:
            vm = si.content.searchIndex.FindByIp(None, ip_address, True)
        else:
            return Response({'message': 'missing arguments'}, status=401)

        print("Found VM: {0}".format(vm.name))
        if format(vm.runtime.powerState) == "poweredOn":
            print("Attempting to power off {0}".format(vm.name))
            task = vm.PowerOffVM_Task()
            print("{0}".format(task.info.state))

        print("Destroying VM from vSphere.")
        task = vm.Destroy_Task()
        print("Done.")
        return Response({'message': 'success'}, status=200)

class GetVMUsage(APIView):
    def printvminfo(self, vm_collect, vm, depth=1):
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
                self.printvminfo(vm_collect, child, depth+1)
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

    def post(self, request):
        
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        
        try:
            si = None
            si = SmartConnectNoSSL(
                host=host,
                user=username,
                pwd=password,
                port=443)
            atexit.register(Disconnect, si)
        except:
            return Response({'message': 'Failure'}, status=401)
        
        if si is None:
            return Response({'message': 'Failure'}, status=401)
        
        content = si.RetrieveContent()

        vm_collect = {}
        for child in content.rootFolder.childEntity:
            if hasattr(child, 'vmFolder'):
                datacenter = child
                vmfolder = datacenter.vmFolder
                vmlist = vmfolder.childEntity
                for vm in vmlist:
                    self.printvminfo(vm_collect, vm)
        # print(vm_collect)
        return Response({'message': 'success', 'host_name': host, 'vm_usage': vm_collect}, status=200)

class setAnnotation(APIView):
    def wait_for_task(self, task):
        ''' wait for a vCenter task to finish '''
        task_done = False
        while not task_done:
            if task.info.state == 'success':
                return task.info.result

            if task.info.state == 'error':
                print("there was an error")
                task_done = True

    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        ip_address = data['ip_address']

        si = None
        try:
            si = SmartConnectNoSSL(host=host,
                                   user=username,
                                   pwd=password,
                                   port=443)
            atexit.register(Disconnect, si)
            content = si.RetrieveContent()
        except IOError as e:
            return Response({'message': 'Failure'}, status=401)
            pass

        if not si:
            return Response({'message': 'Failure'}, status=401)

        vm = si.content.searchIndex.FindByIp(None, ip_address, True)
        if not vm:
            return Response({'message': 'Failure'}, status=401)

        print("Found: {0}".format(vm.name))
        spec = vim.vm.ConfigSpec()
        old_ann = vm.summary.config.annotation
        now = int(round(time.time()*1000))
        now02 = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(now/1000))
        annotation_add = " started at " + now02
        print(old_ann)
        spec.annotation = old_ann.replace("unstarted", annotation_add)

        task = vm.ReconfigVM_Task(spec)
        tasks.wait_for_tasks(si, [task])
        print("Done.")
        return Response({'message': 'success'}, status=200)