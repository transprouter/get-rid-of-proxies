get rid of proxies
==================

Configure once and for all your f proxy.

Testing notes
-------------

```sh
sudo groupadd -r proxy
sudo gpasswd proxy
sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --gid-owner proxy -j REDIRECT --to-ports 3128

sg proxy "go run main.go"
```

socat

    echo -e "GET /resource HTTP/1.1\r\n" | socat TCP4:somehost.tld:80 STDIO
