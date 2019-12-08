from django.db import models
from django.contrib.auth.models import User
from projects.models import *


# from resource.models import Resource


# Create your models here.
class Policy(models.Model):
    DEPLOY_TYPE = (
        ('1', 'K8S'),
        ('2', 'VM')
    )

    ACCOUNT_TYPE = (
        ('1', 'acc_manager'),
        ('2', 'boinc')
    )

    name = models.CharField(max_length=150)
    date_created = models.DateTimeField(blank=True, null=True)
    is_destroy = models.BooleanField(blank=True, null=True, default=True)
    username = models.CharField(max_length=150, blank=True, null=True)
    password = models.CharField(max_length=150, blank=True, null=True)
    deploy_type = models.CharField(max_length=20, choices=DEPLOY_TYPE, default='VM')
    account_type = models.CharField(max_length=20, choices=ACCOUNT_TYPE, default='boinc')
    idle_policy = models.TextField(blank=True, null=True)
    threshold_policy = models.TextField(blank=True, null=True)
    project = models.ForeignKey(Projects, on_delete=models.SET_NULL, null=True)
    user = models.ManyToManyField(User, blank=True)

    #  resource = models.ForeignKey(Resource, on_delete=models.SET_NULL, null=True, related_name="resources")

    class Meta:
        verbose_name = 'Tides Policy'
        verbose_name_plural = 'Tides Policies'

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something
