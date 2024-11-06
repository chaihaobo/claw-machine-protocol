package protocol

import (
	"fmt"
	"net"

	"github.com/chaihaobo/claw-machine-protocol/commands"
)

func (s *Server) listenInternalControl() error {
	listener, err := net.Listen("tcp", s.config.interControlAddr)
	if err != nil {
		return err
	}
	s.internalControlTcpListener = listener
	s.log("listener internal control on %s", s.config.interControlAddr)
	s.acceptInternalControlConnections()
	return nil
}

func (s *Server) acceptInternalControlConnections() {
	for {
		conn, err := s.internalControlTcpListener.Accept()
		if err != nil {
			s.log("accept internal control connection failed: %s", err)
			return
		}
		go s.registerMachineConnection(conn)
	}

}

func (s *Server) registerMachineConnection(conn net.Conn) {
	c := NewConnection(s.config.logger, conn)
	output := &commands.QueryStatusOutput{}
	if err := c.Send(&commands.QueryStatusInput{}, output); err != nil {
		s.log("query machine status failed: %s", err)
		c.Close()
	}
	deviceID := output.DeviceID
	s.connections.Store(deviceID, c)
	c.deviceID = deviceID
	c.statusEventHandler = s.config.statusEventHandler

	if s.registry != nil {
		if err := s.registry.Register(deviceID, s.createDevice(deviceID, conn.RemoteAddr().String())); err != nil {
			s.log("register machine to registry failed:%s", deviceID)
		}
	}
	s.log("register machine successful:%s", deviceID)
}

func (s *Server) createDevice(id, addr string) *Device {
	return NewDevice(id, addr, fmt.Sprintf("%s:%d", localIP(), s.config.externalControlPort))
}

func localIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ip := ipnet.IP
				return ip.String()
			}
		}
	}
	return ""
}

func (s *Server) GetConnection(deviceID string) *Connection {
	if value, ok := s.connections.Load(deviceID); ok {
		return value.(*Connection)
	}
	return nil
}
