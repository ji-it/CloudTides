from django.db import models
from resource.models import *

# Create your models here.
class HostUsage(models.Model):

    date_added = models.DateTimeField(blank=True, null=True)
    host_address = models.TextField()
    host_name = models.TextField(unique=True)
    total_ram = models.FloatField(blank=True, null=True)
    total_cpu = models.FloatField(blank=True, null=True)
    ram_percent = models.FloatField(blank=True, null=True)
    cpu_percent = models.FloatField(blank=True, null=True)
    resource = models.ForeignKey(Resource, on_delete=models.CASCADE)

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something


class VMUsage(models.Model):

    date_added = models.DateTimeField(blank=True, null=True)
    ip_address = models.TextField(unique=True)
    vm_name = models.TextField(blank=True, null=True)
    cpu_usage = models.FloatField(blank=True, null=True)
    mem_usage = models.FloatField(blank=True, null=True)
    resource = models.ForeignKey(Resource, on_delete=models.CASCADE)
    create_time = models.DateTimeField(blank=True, null=True)
    boinc_time = models.DateTimeField(blank=True, null=True)

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something