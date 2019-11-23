from django.conf.urls import url
from django.urls import path, include
from .views import *

app_name = 'users'
urlpatterns = [

    path('', UserListView.as_view(), name='view-all'),
    path('login/', UserLogin.as_view(), name='login'),
    #path('', test, name='test')
]
