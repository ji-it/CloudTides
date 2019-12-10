from django.shortcuts import render
from rest_framework.views import APIView
from .models import *
from rest_framework.response import Response
import datetime
import json
from django.shortcuts import get_object_or_404
from resource.models import *
from django.utils import timezone
import time
from rest_framework.permissions import IsAuthenticated
from rest_framework.authtoken.models import Token


# Create your views here.
class HostPastUsage(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        results = []
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        resources = Resource.objects.filter(user=user)
        for count, resource in enumerate(resources):
            usages = HostUsage.objects.select_related("resource").filter(resource__user=user,
                                                                         resource__host_name=resource.host_name).order_by(
                '-date_added')[:30]
            result = {}
            result[resource.host_name] = {}
            result[resource.host_name]["ram"] = []
            result[resource.host_name]["cpu"] = []
            result[resource.host_name]["time"] = []
            for usage in usages:
                result[resource.host_name]["ram"].insert(0, 100 * usage.ram / usage.resource.total_ram)
                result[resource.host_name]["cpu"].insert(0, 100 * usage.cpu / usage.resource.total_cpu)
                result[resource.host_name]["time"].insert(0, usage.date_added.strftime("%H:%M"))
            results.append(result)
        return Response({'message': 'success', 'results': results, "status": True}, status=200)


class AddHostUsage(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        host_name = data['host_name']
        total_ram = data['total_ram']
        total_cpu = data['total_cpu']
        ram_percent = data['ram_percent']
        cpu_percent = data['cpu_percent']
        date_added = datetime.datetime.now()

        try:
            res = Resource.objects.get(host_address=host_address)
        except:
            return Response({'message': 'resource not registered'}, status=401)

        profile = HostUsage(date_added=date_added, host_address=host_address, host_name=host_name, total_ram=total_ram,
                            total_cpu=total_cpu, ram_percent=ram_percent, cpu_percent=cpu_percent, resource=res)

        try:
            profile.save()
        except:
            return Response({'message': 'object exists'}, status=401)

        return Response({'message': 'success'}, status=200)


class UpdateHostUsage(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        host_name = data['host_name']
        resource = Resource.objects.get(host_name=host_name)
        current_ram = data['current_ram']
        current_cpu = data['current_cpu']
        resource.current_cpu = current_cpu
        resource.current_ram = current_ram
        resource.save()
        date_added = timezone.make_aware(datetime.datetime.now(), timezone.get_default_timezone())
        _ = HostUsage.objects.create(date_added=date_added, ram=current_ram, cpu=current_cpu,
                                     resource=resource)
        return Response({'message': 'host usage recorded'}, status=200)


class DeleteHostUsage(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        host_name = data['host_name']
        try:
            HostUsage.objects.get(host_address=host_address, host_name=host_name).delete()
        except:
            return Response({'message': 'object not found'}, status=401)

        return Response({'message': 'success'}, status=200)


class AddVMUsage(APIView):

    def post(self, request):
        data = json.loads(request.body)
        vms = data["vms"]
        total_vms = data['total_vms']
        for vm in vms.keys():
            ip_address = vm
            cpu_usage = vms[vm]['cpu_usage']
            mem_usage = vms[vm]['memory_usage']
            vm_powered_on = vms[vm]['powered_on']
            host_name = vms[vm]["dc_name"]
            vm_name = vms[vm]['name']
            create_time = vms[vm]['vm_created_time']
            boinc_time = vms[vm]['boinc_start_time']
            direct_host = vms[vm]['direct_host_name']
            vm_total_mem = vms[vm]['max_mem']
            vm_total_cpu = vms[vm]['max_cpu']
            vm_num_cpu = vms[vm]['num_cpu']
            vm_guest_os = vms[vm]['guest_os']

            if boinc_time == 'unstarted':
                boinc_time = None

            date_added = timezone.make_aware(datetime.datetime.now(), timezone.get_default_timezone())

            vm, created = VM.objects.get_or_create(ip_address=ip_address, name=vm_name, direct_host=direct_host)
            if created:
                vm.date_created = create_time
                vm.boinc_time = boinc_time
                vm.total_ram = vm_total_mem
                vm.total_cpu = vm_total_cpu
                vm.num_cpu = vm_num_cpu
                vm.guest_os = vm_guest_os
                resource = Resource.objects.get(datacenter=host_name)
                resource.total_vms = total_vms
                resource.save()
                vm.resource = resource
            vm.current_cpu = cpu_usage
            vm.current_ram = mem_usage
            vm.powered_on = vm_powered_on
            vm.save()
            _ = VMUsage.objects.create(date_added=date_added, mem=mem_usage, cpu=cpu_usage, vm=vm)

        return Response({'message': 'success'}, status=200)


'''
class UpdateVMUsage(APIView):

    def post(self, request):
        data = json.loads(request.body)
        #print(data)
        
        for vm in data.keys():
            ip_address = vm
            vm_name = data[vm]['Name']
            cpu_usage = data[vm]['CPU']
            mem_usage = data[vm]['Memory']
            create_time = data[vm]['CreateTime']
            boinc_time = data[vm]['BOINCTime']
            date_added = datetime.datetime.now()

            obj = get_object_or_404(VMUsage, ip_address=ip_address, vm_name=vm_name)
            obj.cpu_usage = cpu_usage
            obj.mem_usage = mem_usage
            obj.date_added = date_added
            obj.create_time = create_time
            obj.boinc_time = boinc_time

            obj.save(update_fields=['cpu_usage', 'mem_usage', 'date_added', 'create_time', 'boinc_time'])

        return Response({'message': 'VM usage updated'}, status=200)
'''
