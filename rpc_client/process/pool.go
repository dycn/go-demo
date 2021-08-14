package process

import (
	"errors"
	"sync"

	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftPoolConfig struct {
	MaxIdle int //最大空闲,即容量
	Addr    string
}

type thriftPool struct {
	pool chan *thrift.TSocket
	lock sync.RWMutex
}

func NewThriftPool(config *ThriftPoolConfig) (*thriftPool, error) {
	t := &thriftPool{
		pool: make(chan *thrift.TSocket, config.MaxIdle),
	}
	for i := 0; i < config.MaxIdle; i++ {
		socket, err := thrift.NewTSocketConf(config.Addr, &thrift.TConfiguration{})
		//socket, err := thrift.NewTSocket(config.Addr)
		if err != nil {
			t.Free()
			return nil, err
		}
		t.pool <- socket
	}
	return t, nil
}

func (t *thriftPool) Get() (*thrift.TSocket, error) {
	for {
		select {
		case conn := <-t.pool:
			if conn == nil {
				return nil, errors.New("连接池已关闭")
			}
			return conn, nil
			//default:
			//	return nil, errors.New("连接池空了")
		}
	}
}

// Put 放回连接池
func (t *thriftPool) Put(conn *thrift.TSocket) error {
	if conn == nil {
		return errors.New("连接已超时")
	}
	t.lock.Lock()
	if t.pool == nil {
		t.lock.Unlock()
		return (*conn).Close()
	}
	select {
	case t.pool <- conn:
		t.lock.Unlock()
	default:
		return errors.New("连接池已满")
	}
	return nil
}

// Free 清空连接池
func (t *thriftPool) Free() {
	t.lock.Lock()
	t.pool = nil
	t.lock.Unlock()
	return
}
