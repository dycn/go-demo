package process

import (
	"fmt"
	"log"
	"os"
	"zdy_demo/common/rpc"
	"zdy_demo/utility"

	"gopkg.in/ini.v1"
)

func Thrift() {
	config := loadIni("thrift")
	ip, ok := config["ip"]
	if !ok {
		log.Fatalln("ip配置有问题")
	}
	port, ok := config["port"]
	if !ok {
		log.Fatalln("port配置有问题")
	}
	addr := fmt.Sprintf("%s:%s", ip, port)

	protocol, transport := rpc.NewConfig(rpc.Binary, rpc.Buffered)

	//增加测试传输格式的数据
	addStudent(protocol, transport, addr)
	count(protocol, transport, addr)
	startTime := utility.GetCurrentMillisecond()

	shortSocketProcess(protocol, transport, addr)
	longSocketProcess(protocol, transport, addr)
	goroutineSocketProcess(protocol, transport, addr)
	getStudentList(protocol, transport, addr)

	endTime := utility.GetCurrentMillisecond()
	fmt.Printf("总耗时 %d - %d = %d ms", endTime, startTime, endTime-startTime)
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
	//fmt.Println(ret)
	return ret
}
