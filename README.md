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

`integration-tests/docker/` contains a docker-compose file to run the required fake networks and services to run the tests.

* `workstation` hosts 2 user accounts:
  - **direct** who has a direct access to the network
  - **proxied** who has all its TCP traffic redirected to transprouter
* `web.local` hosts a webserver responding to `:80/lost` and `:443/lost`
* `ssh.local` hosts a SSH service listening on `:22`
* `proxy` is a Squid proxy server filtering access to external network
* `web.away` hosts the same webserver as `web.local` but responds with a different content
* `ssh.away` hosts the same SSH service as `ssh.local`
```
                                              |
 +-------------+              +-----------+   |
 | workstation +--------------+ web.local |   |
 +-------------+      |       +-----------+   |
                      |                       |    local
                      |       +-----------+   |    network
                      +-------+ ssh.local |   |
                      |       +-----------+   |
                      |                       |
                      |
                      |
                  +-------+
                  | proxy |
                  +-------+
                      |
                      |
                      |                       |
                      |       +----------+    |
                      +-------+ web.away |    |
                      |       +----------+    |   external
                      |                       |   network
                      |       +----------+    |
                      +-------+ ssh.away |    |
                              +----------+    |
                                              |
```
