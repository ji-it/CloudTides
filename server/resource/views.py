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
from usage.models import *
from django.db.models import Count
from django.views.decorators.cache import never_cache

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
        policy = data['policy']
        platform_type = data['vmtype']
        datacenter_name = data['datacenters']

        policy = Policy.objects.get(pk=policy)
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
                    total_ram = ho.hardware.memorySize / GBFACTOR
                    total_cpu = round(((ho.hardware.cpuInfo.hz / 1e+9) * ho.hardware.cpuInfo.numCpuCores), 0)
                    current_ram = float(ho.summary.quickStats.overallMemoryUsage / 1024.0)
                    current_cpu = float(ho.summary.quickStats.overallCpuUsage / 1024.0)
                    # print(current_cpu, current_ram, total_ram, total_cpu)

                    resource = Resource.objects.get_or_create(host_name=host_name, username=username,
                                                              password=password,
                                                              date_added=date_added,
                                                              host_address=host, platform_type=platform_type,
                                                              total_ram=total_ram,
                                                              total_cpu=total_cpu,
                                                              policy=policy,
                                                              datacenter=datacenter_name,
                                                              current_ram=current_ram, current_cpu=current_cpu,
                                                              is_active=is_active)
                    resource[0].user.set((user,))
                    serialized_obj = ResourceSerializer(resource[0])
                    resources.append(serialized_obj.data)

        return Response({'message': 'success', 'status': True, 'results': resources}, status=200)


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


class DestroyVM(APIView):

    def post(self, request):
        data = json.loads(request.body)
        ip_address = data['ip_address']

        try:
            vm = VM.objects.get(ip_address=ip_address)
            vm.date_destroyed = timezone.make_aware(datetime.datetime.now(), timezone.get_default_timezone())
            vm.is_destroyed = True
            vm.powered_on = False
            vm.save()
        except:
            return Response({'message': 'object not found'}, status=401)

        return Response({'message': 'success'}, status=200)


class ResourceInfo(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        results = []
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        resources = Resource.objects.select_related('policy').filter(user=user)
        serializer = ResourceSerializer(resources, many=True)
        for count, resource in enumerate(resources):
            result = serializer.data[count]
            vms_deployed = VM.objects.select_related("resource").filter(resource__user=user,
                                                                        resource__host_name=resource.host_name)
            result["total_vms"] = vms_deployed.count()
            result["active_vms"] = vms_deployed.filter(is_destroyed=False, powered_on=True).count()
            result["last_deployed"] = vms_deployed.order_by('-date_created')[0:1].get().date_created
            results.append(result)
        return Response({'message': 'success', 'results': results, "status": True}, status=200)


class ResourceVMsInfo(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        results = []
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        resources = Resource.objects.filter(user=user)
        for count, resource in enumerate(resources):
            vms_deployed = VM.objects.select_related("resource").filter(is_destroyed=False, resource__user=user,
                                                                        resource__host_name=resource.host_name)
            result = VMSerializer(vms_deployed, many=True)
            results.append(result.data)
        return Response({'message': 'success', 'results': results, "status": True}, status=200)


class OverviewStats(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        results = {}
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        resources = Resource.objects.filter(user=user)
        total_vms = 0
        for resource in resources:
            total_vms = total_vms + resource.total_vms
        active_tides_vms = VM.objects.filter(is_destroyed=False, powered_on=True).count()

        vm_usage = VMUsage.objects.select_related('vm', 'vm__resource').filter(vm__is_destroyed=False,
                                                                               vm__powered_on=True,
                                                                               vm__resource__user=user)
        power_cpu_usage, power_ram_usage = 0, 0
        for vm in vm_usage:
            power_cpu_usage = power_cpu_usage + (vm.cpu / vm.vm.resource.total_cpu)
            power_ram_usage = power_ram_usage + (vm.mem / vm.vm.resource.total_ram)
        power_ram_usage = power_ram_usage / vm_usage.count()
        power_cpu_usage = power_cpu_usage / vm_usage.count()
        power = 60 * power_cpu_usage + 20 * power_ram_usage

        vm_usage = VMUsage.objects.select_related('vm', 'vm__resource').filter(vm__is_destroyed=True,
                                                                               vm__boinc_time__isnull=False,
                                                                               vm__resource__user=user)
        cost_cpu_usage, cost_ram_usage = 0, 0
        for vm in vm_usage:
            cost_cpu_usage = cost_cpu_usage + (vm.cpu / vm.vm.total_cpu)
            cost_ram_usage = cost_ram_usage + (vm.mem / vm.vm.total_ram)
        cost = 600 * cost_cpu_usage + 200 * cost_ram_usage

        vms_deployed = VM.objects.select_related("resource").filter(resource__user=user)
        total_vms_deployed = vms_deployed.count()
        running_vms_deployed = vms_deployed.filter(is_destroyed=False, powered_on=True).count()
        destroyed_vms_deployed = vms_deployed.filter(is_destroyed=True).count()
        resources_used = vms_deployed.aggregate(Count('resource', distinct=True))

        results["resource"] = {
            'hosts': resources.count(),
            'vms': total_vms,
            'contributing': active_tides_vms
        }
        results["contribution"] = cost
        results["power"] = power
        results["workload"] = {
            'contributed': total_vms_deployed,
            'running': running_vms_deployed,
            'destroyed': destroyed_vms_deployed,
            'hosts_used': resources_used['resource__count']

        }
        return Response({'message': 'success', 'results': results, "status": True}, status=200)
