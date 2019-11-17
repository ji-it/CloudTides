# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.shortcuts import render
from rest_framework.response import Response
from django.contrib.auth.models import Group, User
from rest_framework.views import APIView
from .serializers import *
from .models import vCenter


class UserListView(APIView):

    def get(self, request):
        users = vCenter.objects.all()
        serializer = vCenterSerializer(users, many=True)
        return Response({'status': 'SUCCESS', 'list': serializer.data})
