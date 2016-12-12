package main

import (
	"fmt"
	"net"

	"github.com/jeremiehuchet/go-through-proxies/proxy"
	"github.com/jeremiehuchet/go-through-proxies/xnet"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/go-through-proxies")
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
	fmt.Printf("Connection( %s )\n", connInfo)

	p := new(proxy.DirectProxy)
	p.Forward(connInfo)
}
