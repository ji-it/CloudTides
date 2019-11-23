# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.shortcuts import render
from rest_framework.response import Response
from django.contrib.auth.models import Group, User
from rest_framework.views import APIView
from .serializers import *
from .models import vCenter
from django.contrib.auth import login
from django.contrib.auth import authenticate
from django.http import HttpResponse, JsonResponse
from rest_framework.authtoken.models import Token

class UserListView(APIView):

    def get(self, request):
        users = vCenter.objects.all()
        serializer = vCenterSerializer(users, many=True)
        return Response({'status': 'SUCCESS', 'list': serializer.data})

class UserLogin(APIView):
    
    def post(self, request):
        dic = {}
        username = request.data.get("username")
        password = request.data.get("password")
        user = authenticate(username=username, password=password)
        if user is not None:
            if user.is_active:
                login(request, user)
                token, _ = Token.objects.get_or_create(user=user)
                res = User.objects.get(username=username)
                dic['token'] = token
                dic['priority'] = res.vCenter.priority
                return JsonResponse(dic, status=200)
            else:
                return JsonResponse({'message': 'Unauthorized'}, status=401)
        else:
            return JsonResponse({'message': 'Unauthorized'}, status=401)

def test(request):
    return HttpResponse(status=200)