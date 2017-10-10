# transprouter

Configure once and for all your f proxy.

[![codecov](https://codecov.io/gh/transprouter/transprouter/branch/master/graph/badge.svg)](https://codecov.io/gh/transprouter/transprouter)
[![Build Status](https://travis-ci.org/transprouter/transprouter.svg?branch=master)](https://travis-ci.org/transprouter/transprouter)

## Testing notes

```sh
sudo groupadd -r proxy
sudo gpasswd proxy
sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --gid-owner proxy -j REDIRECT --to-ports 3128

sg proxy "go run main.go"
```

socat

    echo -e "GET /resource HTTP/1.1\r\n" | socat TCP4:somehost.tld:80 STDIO

## integration-tests

Integration tests for transprouter are implemented with cucumber.

`Vagrantfile` at the root of this repository contains a convenient VM configuration with required test tools.
