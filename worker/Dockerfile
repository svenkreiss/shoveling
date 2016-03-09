FROM ubuntu:latest
ENV TERM xterm

RUN apt-get update && apt-get install -y unzip ca-certificates golang-go supervisor curl
ADD https://releases.hashicorp.com/consul/0.6.3/consul_0.6.3_linux_amd64.zip /tmp/consul.zip
RUN cd /tmp \
  && unzip consul.zip \
  && chmod +x consul \
  && mv consul /bin/consul \
  && rm /tmp/consul.zip
RUN mkdir -p /opt/consul
RUN mkdir -p /var/log/consul

ADD . /opt/shoveling
WORKDIR /opt/shoveling
RUN go build -o worker worker.go
ADD supervisord/shoveling.conf /etc/supervisor/conf.d/shoveling.conf
RUN mkdir -p /var/log/shoveling

EXPOSE 8300 8301 8301/udp 8302 8302/udp 8400 8500 8600 8600/udp
CMD supervisord -n -c /etc/supervisor/supervisord.conf
