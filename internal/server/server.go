package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"
	"protocol/pkg/request"
	"protocol/pkg/response"
	"time"
)

const (
	network = "tcp"
)

type Server struct {
	port string
}

func New(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Start(ctx context.Context) error {
	data := make(map[string][]byte)
	// Преобразование сети и порта в TCP-адрес
	tcpAddr, err := net.ResolveTCPAddr(network, s.port)
	if err != nil {
		return err
	}

	// Открытие сокета-прослушивателя
	listener, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Прослушивание порта %s...\\n", s.port)

	for {
		select {
		case <-ctx.Done():
			log.Println("Выход по контексту...")
			return nil
		default:
			// Принятие TCP-соединения от клиента
			listener.SetDeadline(time.Now().Add(1 * time.Second)) // Установка таймаута для Accept
			conn, err := listener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Если таймаут, продолжаем цикл
					continue
				}
				log.Println("Ошибка при принятии соединения:", err)
				continue
			}

			// Обработка запросов клиента в отдельной горутине
			go handleConnection(conn, data, ctx)
		}
	}
}

func handleConnection(conn net.Conn, data map[string][]byte, ctx context.Context) {
	defer conn.Close()

	log.Printf("Подключен клиент %s\\n", conn.RemoteAddr().String())
	defer func() {
		log.Printf("Отключен клиент %s\\n", conn.RemoteAddr().String())
	}()

	connReader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var msgType request.Request
		err := binary.Read(connReader, binary.BigEndian, &msgType)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Ошибка при чтении типа сообщения: %v\n", err)
			return
		}

		switch msgType {
		case request.Info:
			var keyLen int32
			err = binary.Read(connReader, binary.BigEndian, &keyLen)
			if err != nil {
				log.Printf("Ошибка при чтении длины ключа: %v\n", err)
				return
			}

			key := make([]byte, keyLen)
			err = binary.Read(connReader, binary.BigEndian, key)
			if err != nil {
				log.Printf("Ошибка при чтении длины ключа: %v\n", err)
				return
			}

			value, exist := data[string(key)]
			if !exist {
				err = binary.Write(conn, binary.BigEndian, response.OK)
				if err != nil {
					log.Println(err)
				}
				return
			}
			err = binary.Write(conn, binary.BigEndian, response.OK)
			if err != nil {
				log.Println(err)
				return
			}

			err := binary.Write(conn, binary.BigEndian, uint32(len(value)))
			if err != nil {
				log.Println(err)
				return
			}

			err = binary.Write(conn, binary.BigEndian, value)
			if err != nil {
				log.Println(err)
				return
			}
		case request.Record:
			var keyLen uint32
			err = binary.Read(connReader, binary.BigEndian, &keyLen)
			if err != nil {
				log.Printf("Ошибка при чтении длины ключа: %v\n", err)
				return
			}

			key := make([]byte, keyLen)
			err = binary.Read(connReader, binary.BigEndian, key)
			if err != nil {
				log.Printf("Ошибка при чтении длины ключа: %v\n", err)
				return
			}

			var valueLen uint32
			err = binary.Read(connReader, binary.BigEndian, &valueLen)
			if err != nil {
				log.Printf("Ошибка при чтении длины значения: %v\n", err)
				return
			}

			value := make([]byte, valueLen)
			err = binary.Read(connReader, binary.BigEndian, value)
			if err != nil {
				log.Printf("Ошибка при чтении значения: %v\n", err)
				return
			}

			data[string(key)] = value

			err = binary.Write(conn, binary.BigEndian, response.OK)
			if err != nil {
				log.Println(err)
				return
			}
		}

	}

}
