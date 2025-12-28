package client

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"protocol/pkg/request"
	"protocol/pkg/response"
)

func (c *Client) Get(key []byte) ([]byte, error) {
	conn, err := net.DialTCP(network, nil, c.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	connReader := bufio.NewReader(conn)

	keyLenSlice := make([]byte, 0)

	keyLen := binary.BigEndian.AppendUint32(keyLenSlice, uint32(len(key)))

	// Отправка данных в соединение
	_, err = conn.Write([]byte{byte(request.Info)})
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(keyLen)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(key)
	if err != nil {
		return nil, err
	}

	var respCode response.Response
	err = binary.Read(connReader, binary.BigEndian, &respCode)
	if err != nil {
		return nil, err
	}
	switch respCode {
	case response.NotFound:
		return nil, fmt.Errorf("key was not found")
	case response.OK:
		var valueLen uint32
		err = binary.Read(connReader, binary.BigEndian, &valueLen)
		if err != nil {
			return nil, err
		}
		value := make([]byte, valueLen)
		err = binary.Read(connReader, binary.BigEndian, value)
		if err != nil {
			return nil, err
		}

		return value, nil
	case response.Error:
		var errLen int32
		err = binary.Read(connReader, binary.BigEndian, &errLen)

		strErr := make([]byte, errLen)
		err = binary.Read(connReader, binary.BigEndian, strErr)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error of operation: %s", strErr)
	}

	return nil, err
}
