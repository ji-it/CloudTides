from django.shortcuts import render
from rest_framework.response import Response
from rest_framework.views import APIView
import json
import datetime
from .models import Policy

# Create your views here.

class AddPolicy(APIView):

    def post(self, request):
        data = json.loads(request.body)
        name = data['name']
        is_destroy = bool(data['is_destroy'])
        deploy_type = data['deploy_type']
        idle_policy = data['idle_policy']
        date_created = datetime.datetime.now()

        profile = Policy(name=name, date_created=date_created, is_destroy=is_destroy, deploy_type=deploy_type, idle_policy=idle_policy)
        try:
            profile.save()
        except:
            return Response({'message: policy name exists'}, status=401)
        
        return Response({'message': 'success'}, status=200)


class RemovePolicy(APIView):

    def post(self, request):
        data = json.loads(request.body)
        name = data['name']

        try:
            Policy.objects.get(name=name).delete()
        except:
            return Response({'message': 'no matching objects'}, status=401)
        
        return Response({'message': 'success'}, status=200)