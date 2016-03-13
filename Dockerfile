FROM ubuntu:16.04
ENV TERM xterm

RUN apt-get update && \
    apt-get install -y unzip ca-certificates golang supervisor curl git && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y -q mysql-server && \
    apt-get clean

# consul
ADD https://releases.hashicorp.com/consul/0.6.3/consul_0.6.3_linux_amd64.zip /tmp/consul.zip
RUN cd /tmp \
  && unzip consul.zip \
  && chmod +x consul \
  && mv consul /bin/consul \
  && rm /tmp/consul.zip
RUN mkdir -p /opt/consul/services
RUN mkdir -p /opt/consul/data
RUN mkdir -p /var/log/consul

# mysql
RUN mkdir -p /opt/mysql
RUN mkdir -p /var/log/mysql

# shoveling
RUN mkdir -p /opt/go/src/github.com/svenkreiss/shoveling
RUN mkdir -p /opt/go/pkg
RUN mkdir -p /opt/go/bin
ENV GOPATH /opt/go
ADD worker /opt/go/src/github.com/svenkreiss/shoveling/worker
RUN go get github.com/miekg/dns
RUN go install github.com/svenkreiss/shoveling/worker

RUN mkdir -p /opt/shoveling
WORKDIR /opt/shoveling
ADD image/supervisord/shoveling.conf /etc/supervisor/conf.d/shoveling.conf
RUN mkdir -p /var/log/shoveling

EXPOSE 3306 8300 8301 8301/udp 8302 8302/udp 8400 8500 8600 8600/udp
CMD supervisord -n -c /etc/supervisor/supervisord.conf
