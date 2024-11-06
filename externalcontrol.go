package protocol

import (
	"errors"
	"net/http"

	"github.com/chaihaobo/claw-machine-protocol/codec"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidDeviceID = errors.New("invalid device id")
)

const (
	ExternalControlTokenHeader    = "authorization"
	ExternalControlDeviceIDHeader = "device-id"
)

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := s.verifyExternalControlToken(request); err != nil {
		s.log("invalid external control token: %v", err)
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(err.Error()))
		return
	}

	// get the control device id from header
	deviceID := request.Header.Get(ExternalControlDeviceIDHeader)
	if deviceID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(ErrInvalidDeviceID.Error()))
		return
	}
	// verify content type is binary stream
	if request.Header.Get("Content-Type") != "application/octet-stream" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("invalid content type"))
		return
	}

	input, err := codec.Decode(request.Body)
	if err != nil {
		s.log("invalid packet: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	connection := s.GetConnection(deviceID)
	if connection == nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("device not found"))
		return
	}
	var output = codec.EmptyPacket()
	if err := connection.Send(input, output); err != nil {
		s.log("error sending packet: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	rawOutput, err := codec.Encode(output)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(rawOutput)

}

func (s *Server) verifyExternalControlToken(request *http.Request) error {
	if token := request.Header.Get(ExternalControlTokenHeader); s.config.externalControlToken != "" && token != s.config.externalControlToken {
		return ErrInvalidToken
	}
	return nil
}
