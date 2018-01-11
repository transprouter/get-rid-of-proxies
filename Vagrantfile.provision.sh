#!/bin/bash -e

basedir=$(cd $(dirname $0) ; pwd)

apt-get update -qq
apt-get install -y git python-pip nginx openssh-server dnsmasq
git clone git://github.com/mininet/mininet
git -C mininet checkout 2.2.2
./mininet/util/install.sh -nfv
rm -rf mininet
pip install -r $basedir/features/requirements.txt
