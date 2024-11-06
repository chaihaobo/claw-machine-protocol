package protocol

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/chaihaobo/claw-machine-protocol/commands"
)

const (
	defaultInternalControlAddr = ":12268"
	defaultExternalControlPort = 12269
	defaultHeartbeatInterval   = 5 * time.Minute
)

type (
	statusEventHandler func(ctx context.Context, deviceID string, event *commands.ReadStatusOutput)
	serverConfig       struct {
		interControlAddr     string
		externalControlPort  int
		externalControlToken string
		heartbeatInterval    time.Duration
		logger               Logger
		statusEventHandler   statusEventHandler
	}
	Server struct {
		config                     serverConfig
		connections                sync.Map
		registry                   Registry
		internalControlTcpListener net.Listener
		externalControlServer      *http.Server
		heartbeatTicker            *time.Ticker
	}
)

func (s *Server) log(msg string, args ...any) {
	if s.config.logger == nil {
		return
	}
	s.config.logger.Printf(msg+"\n", args...)
}

func (s *Server) Start() error {
	go func() {
		if err := s.listenExternalControl(); err != nil {
			s.log("listen external control failed: %s", err)
		}
	}()
	go s.heartbeat()
	return s.listenInternalControl()
}

func (s *Server) heartbeat() {
	for range s.heartbeatTicker.C {
		s.log("start heartbeat ticker.")
		s.connections.Range(func(deviceID, connection any) bool {
			s.log("heartbeat for device %s", deviceID)
			c := connection.(*Connection)
			var statusOutput commands.QueryStatusOutput
			if err := c.Send(&commands.QueryStatusInput{}, &statusOutput); err != nil {
				s.log("heartbeat failed for device %s: %s", deviceID, err.Error())
				s.connections.Delete(deviceID)
				c.Close()
			}
			if s.registry != nil {
				if err := s.registry.Register(deviceID.(string), s.createDevice(deviceID.(string), c.RemoteAddr())); err != nil {
					s.log("heartbeat register machine to registry failed:%s", deviceID)
				}
			}

			return true
		})
	}
}

func (s *Server) Stop() error {
	ctx := context.Background()
	externalControlShutdownErr := s.externalControlServer.Shutdown(ctx)
	s.heartbeatTicker.Stop()
	// close the connections
	s.closeConnections()
	return multierr.Combine(externalControlShutdownErr, s.internalControlTcpListener.Close())

}

// listenExternalControl 监听外部控制HTTP端口
func (s *Server) listenExternalControl() error {
	s.externalControlServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.externalControlPort),
		Handler: s,
	}
	s.log("listen external control on :%d", s.config.externalControlPort)
	return s.externalControlServer.ListenAndServe()
}

func (s *Server) closeConnections() {
	s.connections.Range(func(deviceID, connection any) bool {
		c := connection.(*Connection)
		if err := c.Close(); err != nil {
			s.log("close device %s connection failed:%s", deviceID, err.Error())
		}
		return true
	})
}

func NewServer(redisClient *redis.Client) *Server {
	config := defaultServerConfig()
	return &Server{
		config:          config,
		registry:        NewRegistry(redisClient),
		heartbeatTicker: time.NewTicker(config.heartbeatInterval),
	}
}

func defaultServerConfig() serverConfig {
	return serverConfig{
		interControlAddr:    defaultInternalControlAddr,
		externalControlPort: defaultExternalControlPort,
		heartbeatInterval:   defaultHeartbeatInterval,
		statusEventHandler: func(ctx context.Context, deviceID string, event *commands.ReadStatusOutput) {
			fmt.Printf("handle device %s status changed to %s \n", deviceID, event.StatusType.Value())
		},
		logger: LoggerFunc(func(msg string, args ...any) {
			if len(args) == 0 {
				fmt.Printf(msg)
				return
			}
			fmt.Printf(msg, args...)
		}),
	}

}
