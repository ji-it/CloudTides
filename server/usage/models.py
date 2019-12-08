from django.db import models
from resource.models import *


# Create your models here.
class HostUsage(models.Model):
    date_added = models.DateTimeField(blank=True, null=True)
    ram = models.FloatField(blank=True, null=True)
    cpu = models.FloatField(blank=True, null=True)
    resource = models.ForeignKey(Resource, on_delete=models.DO_NOTHING, null=True)

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something


class VMUsage(models.Model):
    date_added = models.DateTimeField(blank=True, null=True)
    cpu = models.FloatField(blank=True, null=True)
    mem = models.FloatField(blank=True, null=True)
    vm = models.ForeignKey(VM, on_delete=models.DO_NOTHING, null=True)

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something
