package client

import (
	"net"
)

const (
	network = "tcp"
)

type Client struct {
	addr *net.TCPAddr
}

func NewClient(address string) (Client, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	return Client{tcpAddr}, err
}

func (c *Client) Ping() error {
	// Создание соединения с сервером по TCP-адресу
	conn, err := net.DialTCP(network, nil, c.addr)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}
