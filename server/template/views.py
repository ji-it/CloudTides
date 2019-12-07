from django.shortcuts import render
import os

from django.conf import settings
from rest_framework.response import Response
from rest_framework.views import APIView
import json
import datetime
from .models import *
from django.core.files.storage import default_storage
from django.utils import timezone
from rest_framework.permissions import IsAuthenticated
from rest_framework.authtoken.models import Token
from .serializers import *


# Create your views here.

class AddTemplate(APIView):
    permission_classes = (IsAuthenticated,)

    def post(self, request):

        data = json.loads(request.POST['extra'])
        date_added = timezone.make_aware(datetime.datetime.now(), timezone.get_default_timezone())
        guest_os = data['os']
        template_type = data['source']
        compatibility = data['compat'] and data['compat']
        provisioned_space = data['space'] and data['space']
        memory_size = data['memsize'] and data['memsize']

        name, file = '', ''
        if template_type == "upload":
            name = str(request.FILES['file'])
            file = request.FILES['file']

        template = Template(name=name, date_added=date_added, guest_os=guest_os, compatibility=compatibility, file=file,
                            provisioned_space=provisioned_space, memory_size=memory_size, template_type=template_type)
        try:
            template.save()
        except:
            return Response({'message': 'Template exists'}, status=200)

        return Response({'message': 'success'}, status=200)


class ListTemplate(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        templates = Template.objects.all()
        serializer = TemplateSerializer(templates, many=True)
        return Response({'message': 'success', 'status': True, 'results': serializer.data}, status=200)


class DeleteTemplate(APIView):
    permission_classes = (IsAuthenticated,)

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
