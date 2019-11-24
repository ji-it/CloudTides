from rest_framework import serializers
from .models import *
from rest_framework.validators import UniqueValidator
from django.contrib.auth.models import User


class TidesUserSerializer(serializers.ModelSerializer):
    PRI_CHOICES = (
        ('1', 'Low'),
        ('2', 'Medium'),
        ('3', 'High')
    )
    username = serializers.CharField(
        required=True,
        validators=[UniqueValidator(queryset=User.objects.all())]
    )
    password = serializers.CharField(min_length=4, write_only=True)

    def create(self, validated_data):
        user = User.objects.create_user(validated_data['username'], validated_data['password'])
        return user

    class Meta:
        model = TidesUser
        fields = ["username", "password", "priority"]