from django.db import models
from django.contrib.auth.models import User
from policy.models import *


# Create your models here.
class Resource(models.Model):
    PLATFORM = (
        ('1', 'vsphere'),
        ('2', 'kvm'),
        ('3', 'hyper-v')
    )

    STATUS = (
        ('1', 'idle'),
        ('2', 'busy'),
        ('3', 'contributing')
    )

    # name = models.CharField(max_length=200)
    # id = models.AutoField(primary_key=True)
    date_added = models.DateTimeField(blank=True, null=True)
    host_address = models.TextField(null=True)
    host_name = models.TextField(null=True, unique=True)
    platform_type = models.CharField(max_length=10, choices=PLATFORM, default='vsphere')
    username = models.CharField(max_length=150, null=True)
    password = models.CharField(max_length=128, null=True)
    datacenter = models.CharField(max_length=128, null=True)
    status = models.CharField(max_length=20, choices=STATUS, null=True)
    total_disk = models.FloatField(blank=True, null=True)
    total_ram = models.FloatField(blank=True, null=True)
    total_cpu = models.FloatField(blank=True, null=True)
    current_disk = models.FloatField(blank=True, null=True)
    current_ram = models.FloatField(blank=True, null=True)
    current_cpu = models.FloatField(blank=True, null=True)
    is_active = models.BooleanField(default=False)
    total_jobs = models.IntegerField(blank=True, null=True, default=0)
    job_completed = models.IntegerField(blank=True, null=True, default=0)
    polling_interval = models.IntegerField(blank=True, null=True)
    monitored = models.BooleanField(blank=True, null=True, default=False)
    user = models.ManyToManyField(User, blank=True)
    policy = models.ForeignKey(Policy, on_delete=models.SET_NULL, null=True)

    class Meta:
        verbose_name = 'Tides Resource'
        verbose_name_plural = 'Tides Resources'

    def save(self, *args, **kwargs):
        super().save(*args, **kwargs)
