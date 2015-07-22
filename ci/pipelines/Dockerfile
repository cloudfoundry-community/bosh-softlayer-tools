FROM golang

ENV DEBIAN_FRONTEND noninteractive

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get golang.org/x/tools/cmd/vet

RUN	apt-get update

RUN apt-get -y install apt-utils

RUN	apt-get -y install s3cmd
	
RUN apt-get -y install zlib1g-dev

RUN wget https://opscode-omnibus-packages.s3.amazonaws.com/ubuntu/12.04/x86_64/chefdk_0.6.2-1_amd64.deb && \
	dpkg -i chefdk_0.6.2-1_amd64.deb && \
	rm chefdk_0.6.2-1_amd64.deb

RUN apt-get -y install libgl1-mesa-glx libpng12-0 libpython2.7 libqt4-network libqt4-opengl libqtcore4 libqtgui4 libsdl1.2debian libvpx1 libxcursor1 libxinerama1 libxmu6 psmisc

RUN wget http://download.virtualbox.org/virtualbox/5.0.0/virtualbox-5.0_5.0.0-101573~Ubuntu~trusty_amd64.deb && \
    dpkg -i virtualbox-5.0_5.0.0-101573~Ubuntu~trusty_amd64.deb && \
    rm virtualbox-5.0_5.0.0-101573~Ubuntu~trusty_amd64.deb

RUN apt-get update && \
    apt-get -y install dkms

RUN apt-get -y install vagrant

RUN apt-get install -f
RUN apt-get autoremove
RUN apt-get clean