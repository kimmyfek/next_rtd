FROM alpine:3.6
MAINTAINER Jason Job / Kim Fekete
EXPOSE 3000
RUN mkdir -p /var/next_rtd
COPY nxt-linux-amd64 /usr/local/bin/next_rtd
COPY next/ /usr/local/bin/next
COPY runboard/ /var/next_rtd/runboard/
WORKDIR /usr/local/bin
CMD ./next_rtd --sqlPass=CDFFB0269046B874E09717CF57E6DD43 --sqlUser=rtdro
