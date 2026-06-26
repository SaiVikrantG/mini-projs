package server

import (
	"context"
	"fmt"
	"net"

	"github.com/SaiVikrantG/server/internal/handlers"
)

type Server struct {
	Port     int
	Listener net.Listener
}

func ServerInit(port int) *Server {
	return &Server{
		Port: port,
	}
}

func (s *Server) ServerStart() error {
	conn_port := fmt.Sprintf(":%v", s.Port)
	listener, err := net.Listen("tcp", conn_port)
	if err != nil {
		return err
	}

	s.Listener = listener

	fmt.Printf("Server started at port %v\n", s.Port)
	return nil
}

func (s *Server) ServerListen(ctx context.Context) {
	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Server shutting down")
				return
			default:
				fmt.Println("Temporary error accepting connection:", err)
				continue
			}
		}

		go handlers.HandleConn(conn)
	}
}

func (s *Server) ServerStop(ctx context.Context) error {
	err := s.Listener.Close()

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done(): //not needed as such
		fmt.Println("Shutdown timeout exceeded")
	default:
		fmt.Println("Server stopped gracefully")
	}

	return nil
}
