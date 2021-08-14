package process

import (
	"fmt"
	"log"
	"os"
	"zdy_demo/common/rpc"
	"zdy_demo/common/rpc/thrift/gen-go/students"

	"github.com/apache/thrift/lib/go/thrift"

	"gopkg.in/ini.v1"
)

func initProcessor(className string) thrift.TProcessor {
	//handler
	handler := &ClassMembers{Name: className}

	//processor
	return students.NewClassMemberProcessor(handler)
}

func Thrift() {
	config := loadIni("thrift")
	if !checkConfig(config) {
		os.Exit(1)
	}

	processor := initProcessor(config["className"])

	server := rpc.NewServer(processor, rpc.Binary, rpc.Buffered, fmt.Sprintf("%s:%s", config["ip"], config["port"]))
	err := server.Serve()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func Grpc() {
	//service()
	demo()
}

func loadIni(sectionName string) map[string]string {
	ret := make(map[string]string)
	dir, _ := os.Getwd()
	file := dir + "/conf.d/config.ini"
	config, err := ini.Load(file)
	if err != nil {
		log.Fatalln("配置文件找不到")
	}
	ret["ip"] = config.Section(sectionName).Key("ip").String()
	ret["port"] = config.Section(sectionName).Key("port").String()
	ret["className"] = config.Section(sectionName).Key("classname").String()
	//fmt.Println(ret)
	return ret
}

func checkConfig(config map[string]string) bool {
	_, ok := config["ip"]
	if !ok {
		fmt.Println("ip配置有问题")
		return false
	}
	_, ok = config["port"]
	if !ok {
		fmt.Println("port配置有问题")
		return false
	}
	_, ok = config["className"]
	if !ok {
		fmt.Println("classname配置有问题")
		return false
	}
	return true
}
