package socket

import (
   "atlas-socket/crypto"
   "atlas-socket/request"
   "atlas-socket/session"
   "fmt"
   "log"
   "net"
   "os"
   "time"
)

type Server struct {
   logger         *log.Logger
   sessionService session.Service
   handleSupplier request.Supplier
   ipAddress      string
   port           int
}

func NewServer(l *log.Logger, s session.Service, hs request.Supplier, opts ...ServerOpt) (*Server, error) {
   server := Server{l, s, hs, "0.0.0.0", 5000}
   for _, o := range opts {
      o(&server)
   }
   return &server, nil
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

      sw := &serverWorker{
         logger:         s.logger,
         sessionService: s.sessionService,
         handleSupplier: s.handleSupplier,
      }
      go sw.run(c, sessionId, 4)

      sessionId += 1
   }
}

type serverWorker struct {
   logger         *log.Logger
   sessionService session.Service
   handleSupplier request.Supplier
}

func (sw *serverWorker) run(conn net.Conn, sessionId int, headerSize int) {
   defer conn.Close()

   ses, err := sw.sessionService.Create(sw.logger, sessionId, conn)
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
         sw.handle(sessionId, result)
      }

      header = !header
   }

   sw.logger.Printf("[INFO] Session %d exiting read loop.", sessionId)

   ses.Disconnect()
}

func (sw *serverWorker) handle(sessionId int, p request.Request) {
   go func(sessionId int, reader request.RequestReader) {
      op := reader.ReadUint16()
      h := sw.handleSupplier.Supply(op)
      if h != nil {
         h.Handle(sessionId, reader)
      } else {
         sw.logger.Printf("[INFO] Session %d read a unhandled message with op %05X.", sessionId, op&0xFF)
      }
   }(sessionId, request.NewRequestReader(&p, time.Now().Unix()))
}
