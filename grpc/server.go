package grpc

import (
	"fmt"
	"gome/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
)

var Conf *utils.MeConfig

type GRPC struct {
	Listener net.Listener
	Protocol string
	RPCurl   string
}

func init() {
	confFile, _ := ioutil.ReadFile("config.yaml")
	err := yaml.Unmarshal(confFile, &Conf)
	if err != nil {
		fmt.Println(err)
	}
}

func NewRpcListener() *GRPC {
	host := Conf.GRPCConf.Host
	port := Conf.GRPCConf.Port
	RPCurl := host + ":" + port
	gRPC := &GRPC{Protocol: "tcp", RPCurl: RPCurl}

	var err error
	gRPC.Listener, err = net.Listen(gRPC.Protocol, gRPC.RPCurl)
	if err != nil {
		log.Println("错误:", err)
		panic("服务启动失败")
	} else {
		if Conf.Debug {
			log.Println("服务启动成功，正在监听: " + host + ":" + port)
		}
	}

	return gRPC
}
