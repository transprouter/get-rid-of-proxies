package main

import (
	"fmt"
	"net"

	"github.com/transprouter/transprouter/proxy"
	"github.com/transprouter/transprouter/xnet"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/transprouter")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	viper.SetDefault("pac.url", "http://localhost:80/proxy.pac")

	lnaddr, err := net.ResolveTCPAddr("tcp", ":3128")

	ln, err := net.ListenTCP("tcp", lnaddr)
	if err != nil {
		fmt.Printf("Unable to start server: %s\n", err)
	}
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Printf("Unable to accept incoming connection: %s\n", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {

	connInfo := xnet.Inspect(conn)

	//p := new(proxy.DirectProxy)
	p := proxy.NewHTTPProxy("proxy", 3128)
	p.Forward(connInfo)
}
