package process

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"gopkg.in/ini.v1"
)

func TCPServer() {
	config := loadIni("tcp")
	if !checkTCPConfig(config) {
		os.Exit(1)
	}
	maxClient, err := strconv.Atoi(config["maxClient"])
	if err != nil {
		log.Fatalln(err.Error())

	}
	timeout, err := strconv.Atoi(config["timeOut"])
	if err != nil {
		log.Fatalln(err.Error())
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config["ip"], config["port"]))
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer func() {
		_ = listener.Close()
	}()

	maxGoroutine := make(chan struct{}, maxClient)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		maxGoroutine <- struct{}{}
		go HandleTCPConnection(maxGoroutine, conn, timeout)
	}
}

func UDPServer() {
	config := loadIni("udp")
	if !checkUDPConfig(config) {
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", config["ip"], config["port"]))
	if err != nil {
		log.Fatalln(err.Error())

	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	HandleUDPConnection(conn)
}

// loadIni 加载配置
func loadIni(sectionName string) map[string]string {
	ret := make(map[string]string)
	dir, _ := os.Getwd()
	file := dir + "/conf.d/config.ini"
	config, err := ini.Load(file)
	if err != nil {
		log.Fatalln("配置文件找不到", err.Error())
	}
	ret["ip"] = config.Section(sectionName).Key("ip").String()
	ret["port"] = config.Section(sectionName).Key("port").String()
	ret["timeOut"] = config.Section(sectionName).Key("timeout").String()
	ret["maxClient"] = config.Section(sectionName).Key("maxclient").String()
	return ret
}

// checkTCPConfig 校验配置完整性
func checkTCPConfig(config map[string]string) bool {
	_, ok := config["ip"]
	if !ok {
		fmt.Println("ip配置错误")
		return false
	}
	_, ok = config["port"]
	if !ok {
		fmt.Println("端口配置错误")
		return false
	}
	_, ok = config["timeOut"]
	if !ok {
		fmt.Println("超时配置错误")
		return false
	}
	_, ok = config["maxClient"]
	if !ok {
		fmt.Println("超时配置错误")
		return false
	}
	return true
}

func checkUDPConfig(config map[string]string) bool {
	_, ok := config["ip"]
	if !ok {
		fmt.Println("ip配置错误")
		return false
	}
	_, ok = config["port"]
	if !ok {
		fmt.Println("端口配置错误")
		return false
	}
	return true
}
