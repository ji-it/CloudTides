from django.db import models


# Create your models here.
class Projects(models.Model):
    project_name = models.TextField(null=True)
    has_account_manager = models.BooleanField(default=False)
    url = models.TextField(null=True)

    class Meta:
        verbose_name = 'Tides Project'
        verbose_name_plural = 'Tides Projects'

    def save(self, *args, **kwargs):
        super().save(*args, **kwargs)
