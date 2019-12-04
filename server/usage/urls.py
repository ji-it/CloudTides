from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [
    path('addhost/', AddHostUsage.as_view(), name='add_host_usage'),
    path('updatehost/', UpdateHostUsage.as_view(), name='update_host_usage'),
    path('deletehost/', DeleteHostUsage.as_view(), name='delete_host_usage'),
    path('addvm/', AddVMUsage.as_view(), name='add_vm_usage'),
    #path('updatevm/', UpdateVMUsage.as_view(), name='update_vm_usage'),
    path('deletevm/', DeleteVMUsage.as_view(), name='delete_vm_usage')
]