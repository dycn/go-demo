package process

import (
	"errors"
	"sync"
	"time"
)

type ConnPoolConfig struct {
	MaxIdle int           //最大空闲,即容量
	TimeOut time.Duration //过期时间
	Addr    string
	Factory func() (Conner, error)
}

type Conner interface {
	Close() error
}

type connResource struct {
	conn    *Conner
	timeOut time.Time
}

type connPool struct {
	pool    chan *connResource
	lock    sync.RWMutex
	timeOut time.Duration
	factory func() (Conner, error)
	isClose bool
}

func NewTCPConnPool(config *ConnPoolConfig) (*connPool, error) {
	if config.MaxIdle <= 0 {
		return nil, errors.New("连接池容量设置异常")
	}
	c := &connPool{
		pool:    make(chan *connResource, config.MaxIdle),
		factory: config.Factory,
		timeOut: config.TimeOut,
	}
	for i := 0; i < config.MaxIdle; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Free()
			return nil, err
		}
		connRes := &connResource{
			conn:    &conn,
			timeOut: time.Now().Add(config.TimeOut),
		}
		c.pool <- connRes
	}
	return c, nil
}

func (c *connPool) Get() (*Conner, error) {
	if c.pool == nil {
		return nil, errors.New("连接池已关闭")
	}
	for {
		select {
		case connRes, ok := <-c.pool:
			if !ok {
				return nil, errors.New("连接池已关闭")
			}

			//连接超时
			if connRes.timeOut.Before(time.Now()) {
				//超时关闭
				_ = c.Close(connRes.conn)
				//建立新连接
				conn, err := c.factory()
				if err != nil {
					return nil, err
				}
				connRes := &connResource{
					conn:    &conn,
					timeOut: time.Now().Add(c.timeOut),
				}
				//fmt.Println("超时")
				return connRes.conn, nil
			}

			return connRes.conn, nil
		}
	}

}

// Close 关闭连接
func (c *connPool) Close(conn *Conner) error {
	return (*conn).Close()
}

// Put 放回连接池
func (c *connPool) Put(conn *Conner) error {
	if c.isClose {
		return errors.New("连接池已关闭")
	}
	connRes := &connResource{
		conn:    conn,
		timeOut: time.Now().Add(c.timeOut),
	}
	select {
	case c.pool <- connRes:
	default:
		return errors.New("连接池已满")
	}
	return nil
}

// Free 清空连接池
func (c *connPool) Free() {
	if c.isClose {
		return
	}
	c.lock.Lock()
	c.isClose = true
	close(c.pool)
	for connRes := range c.pool {
		_ = (*(*connRes).conn).Close()
	}
	c.lock.Unlock()
	return
}
