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


# Create your views here.
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
        host_usage = HostUsage.objects.create(date_added=date_added, ram=current_ram, cpu=current_cpu,
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
        # print(data)
        i = 0
        host_address = None
        for vm in data.keys():
            if i == 0:
                host_address = vm
                i += 1
                continue
            ip_address = vm
            vm_name = data[vm]['Name']
            cpu_usage = data[vm]['CPU']
            mem_usage = data[vm]['Memory']
            create_time = data[vm]['CreateTime']
            boinc_time = data[vm]['BOINCTime']
            if boinc_time == 'unstarted':
                boinc_time = None
            date_added = datetime.datetime.now()

            if VMUsage.objects.get(ip_address=ip_address):
                obj = get_object_or_404(VMUsage, ip_address=ip_address, vm_name=vm_name)
                obj.cpu_usage = cpu_usage
                obj.mem_usage = mem_usage
                obj.date_added = date_added
                obj.create_time = create_time
                obj.boinc_time = boinc_time

                obj.save(update_fields=['cpu_usage', 'mem_usage', 'date_added', 'create_time', 'boinc_time'])
                continue

            try:
                res = Resource.objects.get(host_address=host_address)
            except:
                return Response({'message': 'resource not registered'}, status=401)

            profile = VMUsage(ip_address=ip_address, vm_name=vm_name, cpu_usage=cpu_usage, mem_usage=mem_usage,
                              date_added=date_added, resource=res, create_time=create_time, boinc_time=boinc_time)
            try:
                profile.save()
            except:
                return Response({'message': 'object exists'}, status=401)

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


class DeleteVMUsage(APIView):

    def post(self, request):
        data = json.loads(request.body)
        ip_address = data['ip_address']

        try:
            VMUsage.objects.get(ip_address=ip_address).delete()
        except:
            return Response({'message': 'object not found'}, status=401)

        return Response({'message': 'success'}, status=200)
