# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from rest_framework.response import Response
from django.contrib.auth import login, authenticate
from rest_framework.views import APIView
from rest_framework.authtoken.models import Token
from rest_framework import status
from .serializers import *
from .models import Account

from django.conf import settings
from django.core.cache.backends.base import DEFAULT_TIMEOUT
from django.views.decorators.cache import cache_page
from rest_framework.permissions import IsAuthenticated
from rest_framework.authtoken.models import Token

CACHE_TTL = getattr(settings, 'CACHE_TTL', DEFAULT_TIMEOUT)


class UserDetails(APIView):
    permission_classes = (IsAuthenticated,)

    def get(self, request):
        token = request.META.get('HTTP_AUTHORIZATION').split(' ')[1]
        user = Token.objects.get(key=token).user
        profile = Account.objects.get(user=user)
        result = ProfileSerializer(profile).data
        return Response({'message': 'success', 'results': result, "status": True}, status=200)


class UserLogin(APIView):

    def post(self, request):
        json = {}
        username = request.data['username']
        password = request.data['password']
        user = authenticate(username=username, password=password)
        if user is not None:
            if user.is_active:
                login(request, user)
                token, _ = Token.objects.get_or_create(user=user)
                res = User.objects.get(username=username)
                json['userInfo'] = {'priority': res.account.priority, 'username': username}
                json['token'] = token.key
                return Response(json, status=status.HTTP_200_OK)
        return Response({'message': 'Unauthorized'}, status=status.HTTP_401_UNAUTHORIZED)


class UserRegister(APIView):

    def post(self, request):
        serializer = TidesUserSerializer(data=request.data)
        if serializer.is_valid():
            user = serializer.save()
            if user:
                token = Token.objects.create(user=user)
                json = serializer.data
                profile = Account(user=user, priority=request.data['priority'],
                                  company_name=request.data['company_name'])
                try:
                    profile.save()
                except:
                    return Response({'message': 'User already registered'}, status=status.HTTP_200_OK)
                json['token'] = token.key
                profile = Account(user=user, priority=request.data['priority'])
                dic = serializer.data
                dic['priority'] = profile.priority
                json = {'userInfo': dic, 'token': token.key}
                profile.save()
                return Response(json, status=status.HTTP_200_OK)

        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
