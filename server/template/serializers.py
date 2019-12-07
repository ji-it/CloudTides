from rest_framework import serializers
from .models import *
from rest_framework.validators import UniqueValidator


class TemplateSerializer(serializers.ModelSerializer):
    class Meta:
        model = Template
        fields = ["name", "date_added", "guest_os", "compatibility", "provisioned_space",
                  "memory_size", "template_type"]
