FROM golang

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get golang.org/x/tools/cmd/vet

RUN apt-get update && apt-get -y install s3cmd
