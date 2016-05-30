#!/usr/bin/env bash

set -e

go version

base=$( cd "$( dirname "$( dirname "$0" )")"/.. && pwd )

base_gopath=$( cd $base/../../../.. && pwd )

export GOPATH=$base/Godeps/_workspace:$base_gopath:$GOPATH

echo "GOPATH=" $GOPATH


function printStatus {
      if [ $? -eq 0 ]; then
          echo -e "\nSWEET SUITE SUCCESS"
      else
          echo -e "\nSUITE FAILURE"
      fi
  }

trap printStatus EXIT

export BMP_UT_OUTPUT=False

echo -e "\n Cleaning build artifacts..."
go clean

echo -e "\n Formatting packages..."
go fmt ./...

echo -e "\n cd to base of project..."
cd $base

echo -e "\n Unit Testing packages:"
ginkgo -r -p --noisyPendings --skipPackage=integration

echo -e "\n Vetting packages for potential issues..."
go tool vet main config cmds common integration


