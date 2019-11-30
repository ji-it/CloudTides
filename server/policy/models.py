from django.db import models

# Create your models here.
class Policy(models.Model):

    DEPLOY_TYPE = (
        ('1', 'kubernetes'),
        ('2', 'vm')
    )

    name = models.CharField(unique=True, max_length=150)
    date_created = models.DateTimeField(blank=True, null=True)
    is_destroy = models.BooleanField(blank=True, null=True)
    deploy_type = models.CharField(max_length=20, choices=DEPLOY_TYPE, default='vm')
    idle_policy = models.TextField(blank=True, null=True)

    class Meta:
        verbose_name = 'Tides Policy'
        verbose_name_plural = 'Tides Policies'

    def save(self, *args, **kwargs):
        # do something
        super().save(*args, **kwargs)
        # do something