from django.shortcuts import render
from rest_framework.response import Response
from rest_framework.views import APIView
import json
import datetime
from .models import Policy
from django.shortcuts import get_object_or_404
from resource.models import *

# Create your views here.

class AddPolicy(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        host_name = data['host_name']
        name = data['name']
        is_destroy = bool(data['is_destroy'])
        deploy_type = data['deploy_type']
        idle_policy = data['idle_policy']
        date_created = datetime.datetime.now()

        res = get_object_or_404(Resource, host_address=host_address, host_name=host_name)

        profile = Policy(host_address=host_address, host_name=host_name, name=name, date_created=date_created,
                    is_destroy=is_destroy, deploy_type=deploy_type, idle_policy=idle_policy, resource=res)
        try:
            profile.save()
        except:
            return Response({'message: policy name exists'}, status=401)
        
        return Response({'message': 'success'}, status=200)


class UpdatePolicy(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        host_name = data['host_name']
        name = data['name']
        is_destroy = bool(data['is_destroy'])
        deploy_type = data['deploy_type']
        idle_policy = data['idle_policy']

        obj = get_object_or_404(Policy, host_address=host_address, host_name=host_name, name=name)
        obj.is_destroy = is_destroy
        obj.deploy_type = deploy_type
        obj.idle_policy = idle_policy
        obj.save(update_fields=['is_destroy', 'deploy_type', 'idle_policy'])

        return Response({'message': 'success'}, status=200)


class RemovePolicy(APIView):

    def post(self, request):
        data = json.loads(request.body)
        host_address = data['host_address']
        name = data['name']

        try:
            Policy.objects.get(host_address=host_address, name=name).delete()
        except:
            return Response({'message': 'no matching objects'}, status=401)
        
        return Response({'message': 'success'}, status=200)