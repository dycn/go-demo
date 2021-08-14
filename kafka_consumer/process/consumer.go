package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"zdy_demo/module/log_record"

	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
)

var ADDR string
var TOPIC string

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
	topic := config.Section("kafka").Key("topic").String()

	ADDR = fmt.Sprintf("%s:%d", ip, port)
	TOPIC = topic
	//fmt.Println(IP, PORT)
}

func LogRecordConsumer(config *sarama.Config) {
	loadIni()
	groupAc1, err := log_record.NewGroup([]string{ADDR}, "group1", config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = groupAc1.Close()
	}()

	client2, err := log_record.NewClient([]string{ADDR}, config)
	groupAc2, err := sarama.NewConsumerGroupFromClient("group1", client2)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { _ = groupAc2.Close() }()

	client3Slow, err := log_record.NewClient([]string{ADDR}, config)
	groupAc3Slow, err := sarama.NewConsumerGroupFromClient("group1", client3Slow)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { _ = groupAc3Slow.Close() }()

	groupB, err := log_record.NewGroup([]string{ADDR}, "group2", config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { _ = groupB.Close() }()

	ctx := context.Background()
	//写文件
	err = writeTimeFile()
	if err != nil {
		fmt.Println("写起始时间戳出错")
		os.Exit(1)
	}

	go func() {
		err := groupAc1.Consume(ctx, []string{TOPIC}, consumerGroupAClient1Handler{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	go func() {
		err := groupAc2.Consume(ctx, []string{TOPIC}, consumerGroupAClient2Handler{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	go func() {
		err := groupAc3Slow.Consume(ctx, []string{TOPIC}, consumerGroupAClient3SlowHandler{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	err = groupB.Consume(ctx, []string{TOPIC}, consumerGroupBHandler{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	time.Sleep(100 * time.Second)
}

func writeTimeFile() error {
	dir, _ := os.Getwd()
	file := dir + "/start_timestamp.txt"
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	str := strconv.FormatInt(time.Now().UnixNano(), 10)

	_, err = f.WriteString(str + "\n")

	if err != nil {
		log.Println(err.Error())
		return err
	}
	_ = f.Close()
	return nil
	//fmt.Println(time.Now().UnixNano())
}
