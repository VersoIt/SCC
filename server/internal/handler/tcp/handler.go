package tcp

import (
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"io"
	"serverClientClient/internal/model"
	"serverClientClient/internal/service"
	"serverClientClient/pkg/server"
)

const bufferReadSize = 4096

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleConn(conn server.ReadWriteConn) {
	buffer := make([]byte, bufferReadSize)
	data := make([]byte, 0, bufferReadSize)
	for {
		readBytes, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Error("error while reading from TCP-client: ", err)
			return
		}
		data = append(data, buffer[:readBytes]...)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		logrus.Error("error while decoding [base64] data from TCP-client: ", err)
		return
	}
	decodedStr := string(decoded)
	logrus.Infof("received from TCP %s: %s", conn.RemoteAddr(), decodedStr)

	err = h.service.Socket.SaveToDB(model.SocketData{Id: conn.RemoteAddr().String(), Data: decodedStr})
	if err != nil {
		logrus.Error("error while saving data from TCP-client: ", err)
		return
	}
}
