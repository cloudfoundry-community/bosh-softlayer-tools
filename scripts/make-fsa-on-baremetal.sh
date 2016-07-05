#!/bin/bash

set -x

## backup the system settings, as they will be overwritten by stemcell;
cp /boot/grub/menu.lst /boot/grub/menu.lst.bak
cp /etc/fstab /etc/fstab.bak
cp /etc/resolv.conf /etc/resolv.conf.bak
cp /etc/resolv.conf /run/resolvconf/resolv.conf.bak
cp /etc/network/interfaces /etc/network/interfaces.bak

rm -fr /boot/*-generic

## extract the stemcell tgz file to get the image file
#wget -P /tmp http://10.113.109.244/mattcui/bosh-stemcell-${build_num}-softlayer-esxi-ubuntu-trusty-go_agent-raw.tgz

## extract the stemcell image file in /, to overwrite current OS files;
tar -zxvf /tmp/bosh-stemcell-${build_num}-softlayer-esxi-ubuntu-trusty-raw.tgz -C /

## restore the system settings, otherwise this vm couldn't start up on reboot;
cp /boot/grub/menu.lst.bak /boot/grub/menu.lst
cp /etc/fstab.bak /etc/fstab
rm /etc/resolv.conf
cp /run/resolvconf/resolv.conf.bak /run/resolvconf/resolv.conf
ln -sf /run/resolvconf/resolv.conf /etc/resolv.conf
cp /etc/network/interfaces.bak /etc/network/interfaces

apt-get -y --force-yes install grub2
update-grub

echo '-------------Disallowing firstboot.sh to run...-------------'
touch /root/firstboot_done

echo '---------------Cleaning up the build temp files...-------------'
mv /root/first* /tmp
rm -fr /root/*
mv /tmp/first* /root
rm -fr /tmp/*.tgz
rm -fr /fsarchiver
rm -fr /tmp/bm_stemcell

## Archive the entire file system with fsarchiver
mkdir -p /fsarchiver
wget -P /fsarchiver http://10.113.109.244/mattcui/fsarchiver

cd /fsarchiver
chmod +x fsarchiver
./fsarchiver -a -A -e /fsarchiver savefs /fsarchiver/bosh-stemcell-${build_num}-softlayer-baremetal.fsa /dev/sda1 /dev/sda6

## Reset root's password
echo "root:%%bm_root_password%%" | chpasswd

## Reset roots ssh public key
mkdir -p /root/.ssh
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzHIdRPegOIyvJ3zAo+VZ5sjvwQIhwGlLp+Jl/yGRmhxdn5juPWKe8HHXvDPxcN4DGV+th+RD3f4zUKF9JTYejLpYOD4YQkfxbqBXl5EJOAOfkdQb1NSt/HcAi+5/eAmuv3pX3xSv1zv6cHLwAhCqFAOuk5mCDm/UPmXA+FRvC2FWyU3s+p1ctXpFSOjtwgJw6v/P33zOGcbECPJLwlIZsbl1hMH4+6bQAd59AFMow6FZ3ShOmdKa7jBSwFIgZIYDeqMW89SOTSlGySAcB8ykWNyviTum0b+qcSxLHtJpXrljY4k51FzAOGh+1sQfTC3o39Hb+A0IbkNNt3Gx/mezL cuixuex@cn.ibm.com" > /root/.ssh/authorized_keys
chown -R root:root "/root/.ssh"
chmod 0755 "/root/.ssh"
chmod 0600 "/root/.ssh/authorized_keys"
