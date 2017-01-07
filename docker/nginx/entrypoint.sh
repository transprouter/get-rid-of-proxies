#!/bin/sh

set -e

mkdir -p /srv/http
echo -n "$*" > /srv/http/lost

nginx -g 'daemon off;'
