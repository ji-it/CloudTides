from django.conf.urls import url
from django.urls import path, include
from .views import *

urlpatterns = [
    # path('', UserListView.as_view(), name='view-all'),
    path('login/', UserLogin.as_view(), name='login'),
    path('register/', UserRegister.as_view(), name='register'),
    path('get_details/', UserDetails.as_view(), name='get_profile')
]
