from django.db import models
from django.contrib.auth.models import User

# Create your models here.

class vCenter(models.Model):
    PRI_CHOICES = (
        ('1', 'Low'),
        ('2', 'Medium'),
        ('3', 'High')
    )
    # DATABASE FIELDS
    user = models.OneToOneField(User, on_delete=models.CASCADE, blank=True)
    #hostURL = models.URLField(max_length=150)
    #company_name = models.CharField(max_length=300)
    priority = models.CharField(max_length=10, choices=PRI_CHOICES, default='Low')

    # META CLASS
    class Meta:
        verbose_name = 'vCenter'
        verbose_name_plural = 'vCenters'

    # TO STRING METHOD
    def __str__(self):
        return str(self.priority)

    # SAVE METHOD
    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something
