
from setuptools import find_packages, setup

setup(
    name='xac',
    version='0.1.0',
    description='Path management toolkit',
    packages=find_packages(),
    install_requires=[
        "lenses==0.5.0",
        "six",
        "pathlib>=1.0.1"
    ]
)
