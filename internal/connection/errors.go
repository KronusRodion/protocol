package connection

import (
	"encoding/binary"
	"io"

	"github.com/KronusRodion/protocol/pkg/response"
)

func WriteError(w io.Writer, msg string, status response.Response) error {
	_, err := w.Write([]byte{byte(status)})
	if err != nil {
		return err
	}

	valueLen := binary.BigEndian.AppendUint32([]byte(nil), uint32(len(msg)))
	_, err = w.Write(valueLen)
	if err != nil {
		return err
	}

	return err
}
