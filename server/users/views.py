# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from rest_framework.response import Response
from django.contrib.auth import login, authenticate
from rest_framework.views import APIView
from rest_framework.authtoken.models import Token
from rest_framework import status
from .serializers import *
from .models import TidesUser

from django.conf import settings
from django.core.cache.backends.base import DEFAULT_TIMEOUT
from django.views.decorators.cache import cache_page

CACHE_TTL = getattr(settings, 'CACHE_TTL', DEFAULT_TIMEOUT)


class UserListView(APIView):

    def get(self, request):
        users = TidesUser.objects.all()
        serializer = TidesUserSerializer(users, many=True)
        return Response({'message': serializer.data}, status=200)


class UserLogin(APIView):

    def post(self, request):
        json = {}
        username = request.POST.get('username')
        password = request.POST.get('password')
        user = authenticate(username=username, password=password)
        if user is not None:
            if user.is_active:
                login(request, user)
                token, _ = Token.objects.get_or_create(user=user)
                res = User.objects.get(username=username)
                json['token'] = token
                json['priority'] = res.TidesUser.priority
                return Response(json, status=status.HTTP_200_OK)
        return Response({'message': 'Unauthorized'}, status=status.HTTP_400_BAD_REQUEST)


class UserRegister(APIView):

    def post(self, request):
        serializer = TidesUserSerializer(data=request.data)
        if serializer.is_valid():
            user = serializer.save()
            if user:
                token = Token.objects.create(user=user)
                json = serializer.data
                profile = TidesUser(user=user, priority=request.data['priority'],
                                    company_name=request.data['company_name'])
                profile.save()
                json['token'] = token.key
                return Response(json, status=status.HTTP_200_OK)

        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
