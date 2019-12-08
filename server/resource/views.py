# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework import status
from django.core import serializers

from .models import Resource
from pyVim.connect import SmartConnectNoSSL, Disconnect
import atexit
from rest_framework.response import Response
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
from rest_framework.authtoken.models import Token
import json
from .serializers import *
from django.utils import timezone
import pyVmomi
from pyVmomi import vim, vmodl
import datetime
from .utils import *
from django.utils.cache import add_never_cache_headers
from django.contrib.auth.models import User
from django.shortcuts import get_object_or_404
from django.views.decorators.cache import never_cache
import requests
from policy.models import *

GBFACTOR = float(1 << 30)


class ValidateResource(APIView):
    permission_classes = (IsAuthenticated,)

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
        return Response({'message': 'Success', 'results': dic}, status=status.HTTP_200_OK)


class AddResource(APIView):
    permission_classes = (IsAuthenticated,)

    def post(self, request):
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        data = json.loads(request.body)
        host = data['host']
        username = data['uname']
        password = data['password']
        platform_type = data['vmtype']
        datacenter_name = data['datacenters']
        polling_interval = data['polling_interval']

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
            return Response({'message': 'Datacenter not found'}, status=401)

        date_added = timezone.make_aware(datetime.datetime.now(), timezone.get_default_timezone())
        host_address = None
        total_ram = None
        total_cpu = None
        current_ram = None
        current_cpu = None
        is_active = True
        resources = []
        if hasattr(datacenter.hostFolder, 'childEntity'):
            hostFolder = datacenter.hostFolder
            computeResourceList = hostFolder.childEntity
            if len(computeResourceList) == 0:
                return Response({'message': 'No host found'}, status=401)
            for computeResource in computeResourceList:
                hostList = computeResource.host
                for ho in hostList:
                    host_name = ho.name
                    total_ram = ho.hardware.memorySize/GBFACTOR
                    total_cpu = round(((ho.hardware.cpuInfo.hz / 1e+9) * ho.hardware.cpuInfo.numCpuCores), 0)
                    current_ram = float(ho.summary.quickStats.overallMemoryUsage / 1024)
                    current_cpu = float(ho.summary.quickStats.overallCpuUsage / 1024)
                    ram_percent = float(current_ram/total_ram)
                    cpu_percent = float(current_cpu/total_cpu)
                    # print(current_cpu, current_ram, total_ram, total_cpu)
                    #try:
                    resource = Resource.objects.get_or_create(host_name=host_name, username=username,
                                                                  password=password,
                                                                  date_added=date_added,
                                                                  host_address=host, platform_type=platform_type,
                                                                  total_ram=total_ram,
                                                                  total_cpu=total_cpu,
                                                                  datacenter=datacenter_name,
                                                                  current_ram=current_ram, current_cpu=current_cpu,
                                                                  is_active=is_active,
                                                                  polling_interval=polling_interval)
                    data = {}
                    data['host_address'] = host
                    data['host_name'] = host_name
                    #data['total_cpu'] = total_cpu
                    #data['total_ram'] = total_ram
                    data['cpu_percent'] = cpu_percent
                    data['ram_percent'] = ram_percent
                    requests.post("http://127.0.0.1:8000/api/usage/addhost/", data=json.dumps(data))

                    resource[0].user.set((user,))
                    serialized_obj = ResourceSerializer(resource[0])
                    resources.append(serialized_obj.data)
                    #except:
                        #return Response({'message': 'Resource already registered', 'status': False}, status=200)

        return Response({'message': 'success', 'status': True, 'results': resources}, status=200)


class UpdateHost(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        host_name = data['host_name']
        current_cpu = data['current_cpu']
        current_ram = data['current_ram']

        host = get_object_or_404(Resource, host_address=host_address, host_name=host_name)
        host.current_cpu = current_cpu
        host.current_ram = current_ram
        host.save(update_fields=['current_cpu', 'current_ram'])
        
        return Response({'message': 'success'}, status=200)


class ListResource(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        resources = Resource.objects.filter(user=user)
        serializer = ResourceSerializer(resources, many=True)
        return Response({'message': 'success', 'status': True, 'results': serializer.data}, status=200)


class DeleteResource(APIView):
    permission_classes = (IsAuthenticated,)

    def post(self, request):
        data = json.loads(request.body)
        host = data['host']
        username = data['username']
        password = data['password']
        # platform_type = data['platform_type']
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


class AssignPolicy(APIView):
    permission_classes = (IsAuthenticated,)

    def post(self, request):
        data = json.loads(request.body)
        resource_id = data['resource_id']
        policy_id = data['policy_id']

        resource = get_object_or_404(Resource, id=resource_id)
        policy = get_object_or_404(Policy, id=policy_id)

        resource.policy = policy
        resource.save(update_fields=['policy'])

        return Response({'message': 'success'}, status=200)