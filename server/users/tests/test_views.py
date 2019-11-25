"""
    Unit Test file for views
"""
from django.test import TestCase
from django.urls import reverse
from django.contrib.auth.models import User
from rest_framework.authtoken.models import Token
from rest_framework import status


class TidesUserViewTest(TestCase):
    """
    Test View class
    """
    # URL for creating an account.
    create_url = reverse('register')

    @classmethod
    def setUpTestData(cls):
        """
        :return: None
        """
        User.objects.create_user('testuser', 'testpassword')

    def test_create_user(self):
        """
        Ensure we can create a new user and a valid token is created with it.
        """
        data = {
            'username': 'foobar',
            'password': 'somepassword',
            'priority': '3',
            'company_name': 'Test Company'
        }

        response = self.client.post(self.create_url, data, format='json')
        user = User.objects.latest('id')

        # We want to make sure we have two users in the database..
        self.assertEqual(User.objects.count(), 2)
        # And that we're returning a 201 created code.
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # Additionally, we want to return the username upon successful creation.
        self.assertEqual(response.data['username'], data['username'])
        self.assertFalse('password' in response.data)
        token = Token.objects.get(user=user)
        self.assertEqual(response.data['token'], token.key)

    def test_create_user_with_short_password(self):
        """
        Ensure user is not created for password lengths less than 4.
        """
        data = {
            'username': 'foobar',
            'password': 'foo',
            'priority': '3',
            'company_name': 'Test Company'
        }

        response = self.client.post(self.create_url, data, format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(User.objects.count(), 1)
        self.assertEqual(len(response.data['password']), 1)

    def test_create_user_with_no_password(self):
        data = {
            'username': 'foobar',
            'password': '',
            'priority': '3',
            'company_name': 'Test Company'
        }

        response = self.client.post(self.create_url, data, format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(User.objects.count(), 1)
        self.assertEqual(len(response.data['password']), 1)

    def test_create_user_with_no_username(self):
        data = {
            'username': '',
            'password': 'foobar',
            'priority': '3',
            'company_name': 'Test Company'
        }

        response = self.client.post(self.create_url, data, format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(User.objects.count(), 1)
        self.assertEqual(len(response.data['username']), 1)

    def test_create_user_with_preexisting_username(self):
        data = {
            'username': 'testuser',
            'password': 'testuser',
            'priority': '3',
            'company_name': 'Test Company'
        }

        response = self.client.post(self.create_url, data, format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(User.objects.count(), 1)
        self.assertEqual(len(response.data['username']), 1)
