from django.shortcuts import render
from rest_framework.response import Response
from rest_framework.views import APIView
import json
import datetime
from .models import Policy
from django.shortcuts import get_object_or_404
from resource.models import *
from django.contrib.auth.models import User

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
        user_account = data['username']

        res = get_object_or_404(Resource, host_address=host_address, host_name=host_name)
        user = get_object_or_404(User, username=user_account)

        profile = Policy(host_address=host_address, host_name=host_name, name=name, date_created=date_created,
                    is_destroy=is_destroy, deploy_type=deploy_type, idle_policy=idle_policy, resource=res, user=user)
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
        user_account = data['username']
        user = get_object_or_404(User, username=user_account)

        obj = get_object_or_404(Policy, host_address=host_address, host_name=host_name, name=name)
        obj.is_destroy = is_destroy
        obj.deploy_type = deploy_type
        obj.idle_policy = idle_policy
        obj.user = user
        obj.save(update_fields=['is_destroy', 'deploy_type', 'idle_policy', 'user'])

        return Response({'message': 'success'}, status=200)


class ListPolicy(APIView):

    def post(self, request):
        data = json.loads(request.body)
        user_account = data['username']
        user = get_object_or_404(User, username=user_account)
        results = Policy.objects.filter(user=user)
        lis = []
        for result in results:
            dic = {}
            dic['host_address'] = result.host_address
            dic['host_name'] = result.host_name
            dic['policy_name'] = result.name
            dic['is_destroy'] = result.is_destroy
            dic['deploy_type'] = result.deploy_type
            dic['idle_policy'] = result.idle_policy
            lis.append(dic)

        return Response({'message': 'success', 'policies': lis}, status=200)


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