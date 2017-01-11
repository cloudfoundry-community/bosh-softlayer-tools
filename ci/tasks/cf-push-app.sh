#!/usr/bin/env bash
set -e

dir=`dirname "$0"`
source ${dir}/utils.sh

print_title "INSTALL CF CLI..."
curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
mv cf /usr/local/bin
echo "cf version..."
cf --version

print_title "CF PUSH APP..."
sed -i '1 i\nameserver '"${CF-NAMESERVER}"'' /etc/resolv.conf
CF_TRACE=true cf api ${CF-API}
CF_TRACE=true cf login -u ${CF-USERNAME} -p ${CF-PASSWORD}
base=`dirname "$0"`
cf push IICVisit -p ${base}/IICVisit.war
curl iicvisit.${APP-API}/GetEnv|grep "DEA IP"
if [ $? -eq 0 ]; then
   echo "cf push app successful!"
else
   echo "cf push app failed!"
fi