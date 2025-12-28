package client

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/KronusRodion/protocol/pkg/request"
	"github.com/KronusRodion/protocol/pkg/response"
)

func (c *Client) Send(key, value []byte) error {
	conn, err := net.DialTCP(network, nil, c.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	connReader := bufio.NewReader(conn)

	keyLen := binary.BigEndian.AppendUint32([]byte(nil), uint32(len(key)))
	valueLen := binary.BigEndian.AppendUint32([]byte(nil), uint32(len(value)))

	// Отправка данных в соединение
	_, err = conn.Write([]byte{byte(request.Record)})
	if err != nil {
		return err
	}

	_, err = conn.Write(keyLen)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(key))
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(valueLen))
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(value))
	if err != nil {
		return err
	}
	// Чтение данных из соединения
	var respCode response.Response
	err = binary.Read(connReader, binary.BigEndian, &respCode)
	if err != nil {
		return err
	}

	switch respCode {
	case response.OK:
		return nil
	case response.Error:
		strErr, err := connReader.ReadString('\n')
		if err != nil {
			return err
		}
		return fmt.Errorf("error of operation: %s", strErr)
	}
	return nil
}
