from django.db import models
from django.contrib.auth.models import User
from projects.models import Projects
from template.models import *

# from resource.models import Resource


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
    template = models.OneToOneField(Template, on_delete=models.CASCADE)
    deploy_type = models.CharField(max_length=20, choices=DEPLOY_TYPE, default='VM')
    account_type = models.CharField(max_length=20, choices=ACCOUNT_TYPE, default='boinc')
    idle_policy = models.TextField(blank=True, null=True)
    threshold_policy = models.TextField(blank=True, null=True)
    project = models.ForeignKey(Projects, on_delete=models.SET_NULL, null=True)
    user = models.ForeignKey(User, blank=True)

    #  resource = models.ForeignKey(Resource, on_delete=models.SET_NULL, null=True, related_name="resources")

    class Meta:
        verbose_name = 'Tides Policy'
        verbose_name_plural = 'Tides Policies'

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something
