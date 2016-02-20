#!/bin/bash
script_dir="$(dirname "$0")"
. /vagrant/vagrant/provisioners.sh

ensure-dir /var/vagrant

update-apt

install software-properties-common # changed in 14.04
install libpcre3-dev
install libcurl3-openssl-dev

apt-get update

install pkg-config
install git-core
install curl
install golang

config-bash
config-hosts
config-locale

#config-db


chown -R vagrant /vagrant


# init scripts here
cd /vagrant
#sudo -u vagrant ./install
cd -
#sudo -u vagrant ./update


exit 0