package cache_client

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

type CacheClient struct {
	conn net.Conn
	lock sync.Mutex
	id   uint64
}

func NewCacheClient(address string) (*CacheClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &CacheClient{
		conn: conn,
		id:   1,
	}, nil
}

func (c *CacheClient) Close() {
	c.conn.Close()
}

func (c *CacheClient) sendCommand(command string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	id := atomic.AddUint64(&c.id, 1) // Increment the id atomically for each new command
	fullCommand := fmt.Sprintf("%d %s\n", id, command)
	_, err := c.conn.Write([]byte(fullCommand))
	if err != nil {
		return "", err
	}

	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return response, nil
}

func (c *CacheClient) Set(key string, value string) (string, error) {
	return c.sendCommand(fmt.Sprintf("SET %s %s", key, value))
}

func (c *CacheClient) Get(key string) (string, error) {
	return c.sendCommand(fmt.Sprintf("GET %s", key))
}

func (c *CacheClient) Del(key string) (string, error) {
	return c.sendCommand(fmt.Sprintf("DEL %s", key))
}
