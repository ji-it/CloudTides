from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [
    path('validate/', ValidateResource.as_view(), name='validate'),
    path('add/', AddResource.as_view(), name='add'),
    #path('list/', ListResource.as_view(), name='list'),
    path('delete/', DeleteResource.as_view(), name='delete')
]