FROM nginx:stable-alpine
RUN apk add --no-cache --quiet openssl
RUN rm -rf /etc/nginx/conf.d \
 && mkdir -p /etc/nginx/conf.d /etc/nginx/ssl \
 && openssl req -x509 -newkey rsa:4096 \
        -subj "/C=XX/ST=XXXX/L=XXXX/O=XXXX/CN=localhost" \
        -keyout "/etc/nginx/ssl/localhost.key" \
        -out "/etc/nginx/ssl/localhost.crt" \
        -days 3650 -nodes -sha256
ADD listeners.conf /etc/nginx/conf.d
ADD entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
