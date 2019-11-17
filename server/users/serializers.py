from rest_framework.serializers import ModelSerializer, PrimaryKeyRelatedField
from .models import *
from rest_framework import *


class vCenterSerializer(ModelSerializer):
    class Meta:
        model = vCenter
        fields = ["company_name"]
