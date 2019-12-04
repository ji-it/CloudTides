from django.shortcuts import render
import os

from django.conf import settings
from rest_framework.response import Response
from rest_framework.views import APIView
import json
import datetime
from .models import *
from django.core.files.storage import default_storage
# Create your views here.

class AddTemplate(APIView):

    def post(self, request):
        #print(dir(request.POST))
        data = json.loads(request.POST['json'])
        name = str(request.FILES['file'])
        date_added = datetime.datetime.now()
        guest_os = data['guest_os']
        compatibility = data['compatibility']
        provisioned_space = data['provisioned_space']
        memory_size = data['memory_size']
        template_type = data['template_type']
        username = data['username']
        password = data['password']
        
        profile = Template(name=name, date_added=date_added, guest_os=guest_os, compatibility=compatibility,
                    provisioned_space=provisioned_space, memory_size=memory_size, template_type=template_type,
                    username=username, password=password)
        try:
            profile.save()
        except:
            return Response({'message': 'template exists'}, status=401)

        try:
            if request.FILES['file']:
                save_path = os.path.join(settings.MEDIA_ROOT, 'uploads', str(request.FILES['file']))
                path = default_storage.save(save_path, request.FILES['file'])
        except:
            return Response({'message': 'template upload error'}, status=401)
        
        return Response({'message': 'success'}, status=200)


class ListTemplate(APIView):

    def post(self, request):
        file_path = os.path.join(settings.MEDIA_ROOT, 'uploads')
        #for fi in os.listdir(file_path):
            #print(fi)

        return Response(os.listdir(file_path), status=200)



class DeleteTemplate(APIView):

    def post(self, request):

        data = json.loads(request.body)
        name = data['name']

        file_path = os.path.join(settings.MEDIA_ROOT, 'uploads', name)
        try:
            os.remove(file_path)
        except:
            return Response({'message': 'file not exist'}, status=401)

        try:
            Template.objects.get(name=name).delete()
        except:
            return Response({'message': 'no matching objects'}, status=401)

        return Response({'message': 'success'}, status=200)