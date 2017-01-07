FROM alpine
RUN apk add --no-cache --quiet squid
ADD squid.conf /etc/squid/squid.conf
EXPOSE 3128
CMD ["squid", "-N"]
