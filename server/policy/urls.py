from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [
    path('add/', AddPolicy.as_view(), name='add_policy'),
    path('update/', UpdatePolicy.as_view(), name='update_policy'),
    path('remove/', RemovePolicy.as_view(), name='remove_policy')
]