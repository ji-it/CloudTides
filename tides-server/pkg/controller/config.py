import os

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DATABASES = {
    'default': {
        'NAME': 'test',
        'USER': 'postgres',
        'PASSWORD': 'Shen1997',
        # created at the time of password setup
        'HOST': '127.0.0.1',
        'PORT': '5432',
    }
}
FULL_HOSTNAME = "http://localhost:80"