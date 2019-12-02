from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [
    path('add/', AddTemplate.as_view(), name='add_template'),
    path('delete/', DeleteTemplate.as_view(), name='delete_template')
]