package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/showntop/sun-broker/hub"
	"github.com/showntop/sun-broker/server"
)

func main() {
	///解析配置参数
	///多语言、国际化
	srv, err := server.Launch("tcp://localhost:1883")
	if err != nil {
		log.Panic(err)
	}
	log.Infof("server start at port: %s", "1883")
	for {
		conn, err := srv.Accept() //不断的获取新的tcp连接
		if err != nil {
			log.Error(err)
		}
		go hub.Mount(conn)
	}
	defer func() {
		err := srv.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

}
