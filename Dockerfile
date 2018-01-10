FROM debian:stretch
ADD bin/vulcand /usr/local/bin
ADD bin/vctl /usr/local/bin
CMD ["vulcand"]
