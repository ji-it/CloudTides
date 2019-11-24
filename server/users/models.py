from django.db import models
from django.contrib.auth.models import User


# Create your models here.

class TidesUser(models.Model):
    PRI_CHOICES = (
        ('1', 'Low'),
        ('2', 'Medium'),
        ('3', 'High')
    )
    # DATABASE FIELDS
    user = models.OneToOneField(User, on_delete=models.CASCADE, primary_key=True)
    company_name = models.CharField(max_length=200)
    priority = models.CharField(max_length=10, choices=PRI_CHOICES, default='Low')

    # META CLASS
    class Meta:
        verbose_name = 'Tides User'
        verbose_name_plural = 'Tides Users'

    # TO STRING METHOD
    def __str__(self):
        return str(self.priority)

    # SAVE METHOD
    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something