package main

import (
	"fmt"

	"github.com/showntop/sun-broker/hub"
	"github.com/showntop/sun-broker/server"
)

func main() {

	server, err := server.Launch("tcp://localhost:1883")
	if err != nil {
		fmt.Println(err)
		panic("介绍")
	}
	// defer server.close
	fmt.Println("server start listener....")
	for {
		conn, err := server.Accept() //不断的获取新的tcp连接
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(conn)
		go hub.Mount(conn)
	}
}
