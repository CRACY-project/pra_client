package extendedrpc

import (
	"net"
	"net/rpc"

	"github.com/jimbersoftware/pra_client/events"
	"github.com/jimbersoftware/pra_client/logging"
)

type ExtendedServer struct {
	listener net.Listener
}

// NewServer creates a new ExtendedServer
func NewServer() *ExtendedServer {
	return &ExtendedServer{}
}

// Register publishes the receiver's methods in the DefaultServer.
func (s *ExtendedServer) Register(rcvr interface{}) error {
	return rpc.Register(rcvr)
}

// Listen starts the ExtendedServer on the given address.
func (s *ExtendedServer) Listen(network, address string) error {
	var err error
	s.listener, err = net.Listen(network, address)
	return err
}

// Accept starts accepting connections and serving them.
func (s *ExtendedServer) Accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *ExtendedServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		logging.Log(logging.DEBUG, "Handling connection...")
		err := s.serveAndCheck(conn)
		if err != nil {
			globalEvents := events.Get()
			globalEvents.Emit(events.RPC_CLIENT_DISCONNECT, nil)
			logging.Log(logging.WARNING, "Lost connection to external service!", err)
			return
		}
	}
}

func (s *ExtendedServer) Ping(args *PingArgs, reply *PingReply) error {
	// This method does nothing. It's just a way for clients to "ping" the server.
	// respond to the client with a "pong" message.
	return nil
}

func (s *ExtendedServer) serveAndCheck(conn net.Conn) error {
	errChan := make(chan error)
	wrappedConn := &wrapConn{Conn: conn, errChan: errChan}

	go rpc.ServeConn(wrappedConn)

	select {
	case err := <-errChan:
		conn.Close()
		return err
	}
}
