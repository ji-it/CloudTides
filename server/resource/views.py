# -*- coding: utf-8 -*-
from __future__ import unicode_literals

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
from django.contrib.auth.models import User
from django.shortcuts import get_object_or_404

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
        polling_interval = data['polling_interval']
        user_account = data['user_account']

        res = get_object_or_404(User, username=user_account)

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
                    host_name = ho.name
                    total_ram = ho.hardware.memorySize/GBFACTOR
                    total_cpu = round(((ho.hardware.cpuInfo.hz/1e+9)*ho.hardware.cpuInfo.numCpuCores),0)
                    current_ram = float(ho.summary.quickStats.overallMemoryUsage/1024)
                    current_cpu = float(ho.summary.quickStats.overallCpuUsage/1024)
                    #print(current_cpu, current_ram, total_ram, total_cpu)
                    profile = Resource(host_name=host_name, username=username, password=password, date_added=date_added, 
                                host_address=host, platform_type=platform_type, total_ram=total_ram, total_cpu=total_cpu,
                                current_ram=current_ram, current_cpu=current_cpu, is_active=is_active,
                                polling_interval=polling_interval)

                    profile.save()
                    profile.user.add(res)
                    
                    
                    #except:
                        #return Response({'message': 'resource already registered'}, status=401)

        return Response({'message': 'success'}, status=200)


class ListResource(APIView):

    def post(self, request):
        data = json.loads(request.body)
        user_account = data['username']

        user = get_object_or_404(User, username=user_account)
        results = Resource.objects.filter(user=user)
        lis = []
        for result in results:
            dic = {}
            dic['host_address'] = result.host_address
            dic['host_name'] = result.host_name
            lis.append(dic)
        
        return Response({'message': 'success', 'resources': lis}, status=200)



class DeleteResource(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        #platform_type = data['platform_type']
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
                    host_name = ho.name
                    try:
                        Resource.objects.get(host_name=host_name).delete()
                    except:
                        return Response({'message': 'no matching objects'}, status=401)
        
        return Response({'message': 'success'}, status=200)