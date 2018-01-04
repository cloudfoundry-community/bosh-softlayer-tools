#!/usr/bin/env bash
set -e

: ${SYSTEM_DOMAIN:?}

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

cf_artifacts_path=$(find ${pwd} -name "cf-artifacts*.tgz" | head -n 1)
tar -zxvf ${cf_artifacts_path} -C ${deployment_dir}

cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Installing Cloud Foundry Client.\033[0m"
# apt-get install wget -y
apt-get install apt-transport-https -y
wget -q -O - https://packages.cloudfoundry.org/debian/cli.cloudfoundry.org.key | apt-key add -
echo "deb http://packages.cloudfoundry.org/debian stable main" | tee /etc/apt/sources.list.d/cloudfoundry-cli.list
apt-get update
apt-get install cf-cli -y

echo -e "\n\033[32m[INFO] Using cf $(cf --version).\033[0m"

echo -e "\n\033[32m[INFO] Logging in CF.\033[0m"
name_server=$(cat deployment/director-hosts | awk '{print $1}')
cat /etc/resolv.conf
sed -i '1 i\nameserver '"${name_server}"'' /etc/resolv.conf
CF_API=api.${SYSTEM_DOMAIN}

CF_TRACE=true cf api ${CF_API} --skip-ssl-validation
echo 1 | CF_TRACE=true cf login -u admin -p $(bosh-cli int ${deployment_dir}/cf-creds.yml --path /cf_admin_password)
cf target -o "cpi_ng" || (cf create-org "cpi_ng" && cf target -o "cpi_ng")
cf target -s "dev" || (cf create-space "dev" && cf target -s "dev")

echo -e "\n\033[32m[INFO] Listing apps in CF before pushing.\033[0m"
cf apps

echo -e "\n\033[32m[INFO] Pushing CF app.\033[0m"
git clone https://github.com/cloudfoundry-samples/cf-sample-app-nodejs.git
pushd cf-sample-app-nodejs
    sed -i '$d' manifest.yml
    echo "  routes:" >> manifest.yml
    echo "    - route: cf-nodejs-sample.${SYSTEM_DOMAIN}" >> manifest.yml
    cf push

    echo -e "\n\033[32m[INFO] Listing apps in CF after pushing.\033[0m"
    cf apps

    echo -e "\n\033[32m[INFO] Verifying sample app.\033[0m"
    response=$(curl --write-out %{http_code} --silent --output /dev/null cf-nodejs-sample.${SYSTEM_DOMAIN})
    if [[ "$response" == "200" ]]; then
        echo -e "\n\033[32m[INFO] The Node.js sample app is executed normally.\033[0m"
        curl cf-nodejs-sample.${SYSTEM_DOMAIN}
    else
        echo -e "\n\033[31m[ERROR] The Node.js sample app is not executed normally.\033[0m"
        exit 1
    fi
popd

echo -e "\n\033[32m[INFO] Saving cf artifacts.\033[0m"
pushd cf-artifacts
    tar -zcvf /tmp/cf_artifacts.tgz ./ >/dev/null 2>&1
popd
mv /tmp/cf_artifacts.tgz ./cf-artifacts
