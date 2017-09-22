#!/usr/bin/env bash
set -e

: ${CF_API:?}

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

tar -zxvf cf-artifacts-comm/cf-artifacts-comm.tgz -C ${deployment_dir}


echo -e "\n\033[32m[INFO] Installing Cloud Foundry Client.\033[0m"
apt-get install wget -y
wget -q -O - https://packages.cloudfoundry.org/debian/cli.cloudfoundry.org.key | apt-key add -
echo "deb http://packages.cloudfoundry.org/debian stable main" | tee /etc/apt/sources.list.d/cloudfoundry-cli.list
apt-get update
apt-get install cf-cli -y
echo -e "\n\033[32m[INFO] Using cf $(cf --version).\033[0m"

echo -e "\n\033[32m[INFO] Pushing CF app.\033[0m"
name_server=$(cat deployment/director-hosts | awk '{print $1}')

cat /etc/resolv.conf
sed -i '1 i\nameserver '"${name_server}"'' /etc/resolv.conf

CF_TRACE=true cf api ${CF_API} --skip-ssl-validation
CF_TRACE=true cf login -u ${CF_USERNAME} -p ${CF_PASSWORD}
cf login -u admin -p $(bosh2 int ${deployment_dir}/cf-vars.yml --path /cf_admin_password)
cf create-org cpi_ng
cf target -o "cpi_ng"
cf create-space dev
cf target -s "dev"
cf apps
git clone https://github.com/cloudfoundry-samples/cf-sample-app-nodejs.git
cd cf-sample-app-nodejs/
cf push
cf apps
response=$(curl --write-out %{http_code} --silent --output /dev/null cf-nodejs-sample.${system-domain})
if [[ "$response" == "200" ]]; then
	echo -e "\n\033[32m[INFO] Node.js sample app is executed normally.\033[0m"
else
	echo -e "\n\033[31m[ERROR] Node.js sample app is not executed normally.\033[0m"
	exit 1
fi

echo -e "\n\033[32m[INFO] Saving cf artifacts.\033[0m"
pushd cf-artifacts
  tar -zcvf /tmp/cf_artifacts.tgz ./ >/dev/null 2>&1
popd
mv /tmp/cf_artifacts.tgz ./cf-artifacts
