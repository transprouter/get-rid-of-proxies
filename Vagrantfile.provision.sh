#!/bin/bash -e

apt-get update
apt-get install -y mininet python-pip nginx
pip install behave PyHamcrest
