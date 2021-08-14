package main

import (
	"fmt"
	"os"
	"zdy_demo/kafka_producer/process"
	"zdy_demo/module/log_record"

	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
)

var ADDR string
var TOPIC string

func main() {
	loadIni()

	config := log_record.NewConfig()
	client, err := sarama.NewSyncProducer([]string{ADDR}, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		_ = client.Close()
	}()

	for i := 0; i < 300000; i++ {
		process.OnProcessLog(client, TOPIC)
	}
}

func loadIni() {
	dir, _ := os.Getwd()
	file := dir + "/conf.d/config.ini"
	config, err := ini.Load(file)
	if err != nil {
		fmt.Println("配置文件找不到")
		os.Exit(1)
	}
	ip := config.Section("kafka").Key("ip").String()
	port, err := config.Section("kafka").Key("port").Int64()
	if err != nil {
		fmt.Println("端口配置有问题")
		os.Exit(1)
	}
	TOPIC = config.Section("kafka").Key("topic").String()

	ADDR = fmt.Sprintf("%s:%d", ip, port)
	//fmt.Println(IP, PORT)
}
