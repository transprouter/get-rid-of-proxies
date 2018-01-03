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

## functional tests

Functional tests are implemented with cucumber (behave) and mininet.

A `Vagrantfile` at the root of this repository contains a convenient VM configuration with required test tools.

**cheatsheet**

- mininet requires root privileges
- `sudo python testnet/topology.py` drops you into a shell to run mininet commands and test the network topology
- `sudo behave basic.feature` runs a cucumber feature file
