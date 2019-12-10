from rest_framework import serializers
from .models import *
from rest_framework.validators import UniqueValidator


class ResourceSerializer(serializers.ModelSerializer):
    cpu_percent = serializers.SerializerMethodField()
    ram_percent = serializers.SerializerMethodField()
    policy_name = serializers.SerializerMethodField()

    class Meta:
        model = Resource
        fields = ["id", "date_added", "host_name", "status", "policy_name",
                  "platform_type", "datacenter", "total_cpu", "total_ram", "total_disk",
                  "current_ram", "current_cpu", "is_active", "total_jobs", "ram_percent",
                  "job_completed", "monitored", "cpu_percent"]

    def get_cpu_percent(self, obj):
        return obj.current_cpu / obj.total_cpu

    def get_ram_percent(self, obj):
        return obj.current_ram / obj.total_ram

    def get_policy_name(self, obj):
        return obj.policy.name


class VMSerializer(serializers.ModelSerializer):
    cpu_percent = serializers.SerializerMethodField()
    ram_percent = serializers.SerializerMethodField()

    class Meta:
        model = VM
        fields = ["id", "date_created", "date_destroyed", "boinc_time", "ip_address", "name",
                  "total_cpu", "total_ram", "total_disk", "current_ram", "current_cpu", "powered_on", "guest_os",
                  "ram_percent", "cpu_percent"]

    def get_cpu_percent(self, obj):
        return obj.current_cpu / obj.total_cpu

    def get_ram_percent(self, obj):
        return obj.current_ram / obj.total_ram
