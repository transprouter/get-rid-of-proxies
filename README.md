# integration-tests

Integration tests for transprouter implemented with cucumber.

`docker/` contains a docker-compose file to run the required fake networks and services to run the tests.

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
