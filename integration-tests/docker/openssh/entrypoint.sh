#!/bin/sh

set -e

ssh-keygen -A

exec "$@"
