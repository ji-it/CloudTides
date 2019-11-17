from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [

    url(r'^api/users/', UserListView.as_view(), name='view-all'),

]
