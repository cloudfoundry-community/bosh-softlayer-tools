FROM golang

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get golang.org/x/tools/cmd/vet

RUN apt-get update && \
	apt-get -y install apt-utils
	
RUN	apt-get update && \
	apt-get -y install s3cmd

RUN wget https://opscode-omnibus-packages.s3.amazonaws.com/ubuntu/12.04/x86_64/chefdk_0.6.2-1_amd64.deb && \
	dpkg -i chefdk_0.6.2-1_amd64.deb && \
	rm chefdk_0.6.2-1_amd64.deb

RUN apt-get update && \
	apt-get -y install vagrant

RUN apt-get autoremove
RUN apt-get clean