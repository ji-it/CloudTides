from rest_framework import serializers
from .models import *
from rest_framework.validators import UniqueValidator


class PolicySerializer(serializers.ModelSerializer):
    project_name = serializers.SerializerMethodField()

    class Meta:
        model = Policy
        fields = ["id", "date_created", "name", "is_destroy", "deploy_type",
                  "idle_policy", "threshold_policy", "project_name"]

    def get_project_name(self, obj):
        return obj.project.project_name
