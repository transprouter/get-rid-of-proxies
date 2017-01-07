FROM alpine
RUN apk add --no-cache --quiet dnsmasq \
 && echo "user=root" >> /etc/dnsmasq.conf
EXPOSE 53 53/udp
ENTRYPOINT ["dnsmasq", "-k"]
