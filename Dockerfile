FROM alpine:3.6
MAINTAINER Jason Job / Kim Fekete
EXPOSE 3000
RUN mkdir -p /var/next_rtd
COPY next_rtd /usr/local/bin/next_rtd
COPY runboard/ /var/next_rtd/runboard/
CMD next_rtd --parse=true --dbDir /var/next_rtd/next.db --sourceDir /var/next_rtd/runboard/
