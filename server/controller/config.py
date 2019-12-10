import os

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DATABASES = {
    'default': {
        'NAME': 'tides2',
        'USER': 'postgres',
        'PASSWORD': 'password',
        'HOST':  'localhost',
        'PORT': '5432',
    }
}
