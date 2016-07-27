#!/usr/bin/env bash

set -e -x

export BOSHIO_BEARER_TOKEN=$boshio_bearer_token

cd light-stemcell-final
export STEMCELL=`ls light*.tgz`

sha1=`sha1sum ${STEMCELL}`
sha1s=$(echo $sha1 | tr " " "\n")
for eachline in $sha1s; do
    if [ -n "${eachline}" ]; then
    	echo "The checksum (sha1) for ${STEMCELL} is ${eachline}, uploading to bosh.io..."
        for((i=1;i<=3;i++));do
        	# curl -X POST 'https://bosh.io/checksums/$BAT_STEMCELL' -d 'sha1=${eachline}' -H 'Authorization: bearer $BOSHIO_BEARER_TOKEN' && exit 0
        	echo "curl -X POST 'https://bosh.io/checksums/$STEMCELL' -d 'sha1=${eachline}' -H 'Authorization: bearer ${BOSHIO_BEARER_TOKEN}'"
        done
        echo "Failed to notify bosh.io with checksum ${eachline} after retrying 3 times"
        exit 1
    else
        echo "Fail to generate checksum for the stemcell"
        exit 1
    fi
    break
done