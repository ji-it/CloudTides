import os

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DATABASES = {
    'default': {
        'NAME': 'Test4',
        'USER': 'postgres',
        'PASSWORD': 't6bB2T5KoQuPq6DrpWxJa3rYKVjIpOCtVSrKyBMB8PHcMShkidcQo8Kjn1lcXswB',
        # created at the time of password setup
        'HOST': '10.11.16.83',
        'PORT': '30123',
    }
}
