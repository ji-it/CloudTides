from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [
    path('validate/', ValidateResource.as_view(), name='validate'),
    path('add/', AddResource.as_view(), name='add'),
    path('delete_host/', DeleteResource.as_view(), name='delete'),
    path('list/', ListResource.as_view(), name='list'),
    path('update/', UpdateHost.as_view(), name='update'),
    path('assign_policy/', AssignPolicy.as_view(), name='assign_policy'),
    path('destroy_vm/', DestroyVM.as_view(), name='destroy_vm'),
    path('overview/', OverviewStats.as_view(), name='overview_stats'),
    path('get_details/', ResourceInfo.as_view(), name='resource_info'),
    path('get_vm_details/', ResourceVMsInfo.as_view(), name='vm_info'),
    path('toggle_active/', ToggleActive.as_view(), name='toggle_active')
]
