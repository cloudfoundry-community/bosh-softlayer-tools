FROM golang

ENV DEBIAN_FRONTEND noninteractive

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get golang.org/x/tools/cmd/vet

RUN apt-get update && \
    apt-get -y install apt-utils

# BEGIN stemcell requirements	
RUN apt-get -y install zlib1g-dev

RUN wget https://opscode-omnibus-packages.s3.amazonaws.com/ubuntu/12.04/x86_64/chefdk_0.9.0-1_amd64.deb && \
	dpkg -i chefdk_*.deb && \
	rm chefdk_*.deb

RUN apt-get -y install libdevmapper1.02.1 libgl1-mesa-glx libpython2.7 libqt4-network libqt4-opengl libqtcore4 libqtgui4 libsdl1.2debian libvpx1 libx11-6 libxcursor1 libxinerama1 libxmu6 psmisc

RUN wget http://download.virtualbox.org/virtualbox/5.0.0/virtualbox-5.0_5.0.0-101573~Ubuntu~trusty_amd64.deb && \
    dpkg -i virtualbox-5.0_5.0.0-101573~Ubuntu~trusty_amd64.deb && \
    rm virtualbox-5.0_5.0.0-101573~Ubuntu~trusty_amd64.deb

RUN apt-get -y install ruby-dev

RUN apt-get -y install vagrant

RUN apt-get -y install python-pip
RUN pip install s3cmd
RUN pip install python-swiftclient

RUN apt-get -y install rsync

# BEGIN ec2 api tools
RUN apt-get -y install unzip

RUN wget http://s3.amazonaws.com/ec2-downloads/ec2-api-tools.zip && \
    mkdir /usr/local/ec2 && \
    unzip ec2-api-tools.zip -d /usr/local/ec2 && \
    rm ec2-api-tools.zip

RUN apt-get -y install openjdk-7-jre

ENV JAVA_HOME /usr/lib/jvm/java-7-openjdk-amd64/jre
ENV EC2_HOME /usr/local/ec2/ec2-api-tools-1.7.5.1
ENV PATH $PATH:$EC2_HOME/bin

RUN apt-get install -f
RUN apt-get autoremove
RUN apt-get clean
RUN apt-get -y upgrade
RUN apt-get -y update

#Install ruby and gems to install bosh CLI
RUN apt-get -y install zlib1g-dev
RUN apt-get -y install ruby ruby-dev
RUN apt-get -y install rubygems
RUN apt-get -y install build-essential
RUN apt-get -y install sqlite3 libsqlite3-dev
RUN gem install bosh_cli bosh_cli_plugin_micro --no-ri --no-rdoc

#Install slcli
RUN pip install softlayer
