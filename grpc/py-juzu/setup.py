#!/usr/bin/env python
from setuptools import setup

setup(
    name='cobalt-juzu',
    python_requires='>=3.5.0',
    version='0.10.0',
    description='This client library is designed to support the Cobalt API for speech diarization with Juzu',
    author='Cobalt Speech and Language Inc.',
    maintainer_email='tech@cobaltspeech.com',
    url='https://cobaltspeech.github.io/sdk-juzu',
    packages=["juzu"],
    install_requires=[
        'googleapis-common-protos==1.53.0',
        'grpcio-tools==1.42.0',
        'grpcio==1.42.0',
    ]
)
