from django.db import models
import os
from django.conf import settings


# Create your models here.

class Template(models.Model):
    TEMPLATE_TYPE = (
        ('1', 'datastore'),
        ('2', 'upload')
    )

    name = models.CharField(max_length=150)
    date_added = models.DateTimeField(blank=True, null=True)
    guest_os = models.CharField(max_length=100)
    compatibility = models.CharField(max_length=100, blank=True, null=True)
    provisioned_space = models.FloatField(blank=True, null=True)
    memory_size = models.FloatField(blank=True, null=True)
    template_type = models.CharField(max_length=20, choices=TEMPLATE_TYPE, null=True)
    username = models.TextField(null=True)
    password = models.TextField(null=True)

    class Meta:
        verbose_name = 'Tides Template'
        verbose_name_plural = 'Tides Templates'

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something
