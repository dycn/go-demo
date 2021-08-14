package rpc

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type ProtocalType uint8
type TransportType uint8

const (
	Binary ProtocalType = iota
	Compact
	JSON
	SimpleJSON
)

const (
	Block TransportType = iota
	Buffered
)

func NewServer(processor thrift.TProcessor, protocolType ProtocalType, transportType TransportType, addr string) *thrift.TSimpleServer {
	protocolFactory := getTProtocol(protocolType)
	transportFactory := getTransport(transportType)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
}

func NewConfig(protocolType ProtocalType, transportType TransportType) (thrift.TProtocolFactory, thrift.TTransportFactory) {
	protocolFactory := getTProtocol(protocolType)
	transportFactory := getTransport(transportType)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	return protocolFactory, transportFactory
}

func getTProtocol(protocol ProtocalType) thrift.TProtocolFactory {
	switch protocol {
	case Compact:
		return thrift.NewTCompactProtocolFactoryConf(&thrift.TConfiguration{})
		//return thrift.NewTCompactProtocolFactory()
	case JSON:
		return thrift.NewTJSONProtocolFactory()
	case SimpleJSON:
		return thrift.NewTSimpleJSONProtocolFactoryConf(&thrift.TConfiguration{})
		//return thrift.NewTSimpleJSONProtocolFactory()
	case Binary:
		fallthrough
	default:
		return thrift.NewTBinaryProtocolFactoryConf(&thrift.TConfiguration{})
		//return thrift.NewTBinaryProtocolFactoryDefault()
	}
}

func getTransport(transport TransportType) thrift.TTransportFactory {
	var transportFactory thrift.TTransportFactory
	switch transport {
	case Block:
		transportFactory = thrift.NewTTransportFactory()
	case Buffered:
		fallthrough
	default:
		transportFactory = thrift.NewTBufferedTransportFactory(4096)
	}
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	return transportFactory
}
