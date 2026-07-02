package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	apperrors "github.com/SaiVikrantG/server/internal/errors"
	"github.com/SaiVikrantG/server/internal/handlers"
	"github.com/SaiVikrantG/server/internal/parser"
	"github.com/SaiVikrantG/server/internal/response"
)

type Server struct {
	Port     int
	Listener net.Listener
	Handler  handlers.Handler
}

func ServerInit(port int, h handlers.Handler) *Server {
	return &Server{
		Port:    port,
		Handler: h,
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

func processRequest(conn net.Conn, h handlers.Handler) {
	defer conn.Close()

	w := response.NewResponseWriter(conn)

	req, err := parser.Parse(conn)
	if err != nil {
		var httpErr *apperrors.HTTPError
		if errors.As(err, &httpErr) {
			w.Write(httpErr.StatusCode, []byte(httpErr.Message))
			return
		}
		return
	}

	h.ServeHTTP(w, req)
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

		go processRequest(conn, s.Handler)
	}
}

func (s *Server) ServerStop(ctx context.Context) error {
	err := s.Listener.Close()

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		fmt.Println("Shutdown timeout exceeded")
	default:
		fmt.Println("Server stopped gracefully")
	}

	return nil
}
