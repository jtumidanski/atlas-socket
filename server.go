package socket

import (
	"fmt"
	"github.com/jtumidanski/atlas-socket/crypto"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/jtumidanski/atlas-socket/session"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	logger         *log.Logger
	sessionService session.Service
	ipAddress      string
	port           int
	handlers       map[uint16]request.Handler
}

func NewServer(l *log.Logger, s session.Service, opts ...ServerOpt) (*Server, error) {
	server := Server{
		l,
		s,
		"0.0.0.0",
		5000,
		make(map[uint16]request.Handler),
	}
	for _, o := range opts {
		o(&server)
	}
	return &server, nil
}

func (s *Server) RegisterHandler(op uint16, handler request.Handler) {
	s.handlers[op] = handler
}

func (s *Server) Run() {
	s.logger.Printf("[INFO] Starting tcp server on %s:%d", s.ipAddress, s.port)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ipAddress, s.port))
	if err != nil {
		s.logger.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer lis.Close()

	sessionId := 0

	for {
		c, err := lis.Accept()
		if err != nil {
			s.logger.Println("Error connecting:", err.Error())
			return
		}

		s.logger.Println("[INFO] Client " + c.RemoteAddr().String() + " connected.")

		go s.run(c, sessionId, 4)

		sessionId += 1
	}
}

func (s *Server) run(conn net.Conn, sessionId int, headerSize int) {
	defer conn.Close()

	ses, err := s.sessionService.Create(s.logger, sessionId, conn)
	if err != nil {
		return
	}

	header := true
	readSize := headerSize

	for {
		buffer := make([]byte, readSize)

		if _, err := conn.Read(buffer); err != nil {
			break
		}

		if header {
			readSize = crypto.PacketLength(buffer)
		} else {
			readSize = headerSize

			result := buffer
			if ses.ReceiveAESOFB() != nil {
				result = ses.ReceiveAESOFB().Decrypt(buffer, true, true)
			}
			s.handle(sessionId, result)
		}

		header = !header
	}

	s.logger.Printf("[INFO] Session %d exiting read loop.", sessionId)

	ses.Disconnect()
}

func (s *Server) handle(sessionId int, p request.Request) {
	go func(sessionId int, reader request.RequestReader) {
		op := reader.ReadUint16()
		if h, ok := s.handlers[op]; ok {
			h(sessionId, reader)
		} else {
			s.logger.Printf("[INFO] Session %d read a unhandled message with op %05X.", sessionId, op&0xFF)
		}
	}(sessionId, request.NewRequestReader(&p, time.Now().Unix()))
}
