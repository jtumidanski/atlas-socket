package handlers

import (
   "github.com/jtumidanski/atlas-socket/request"
   "github.com/jtumidanski/atlas-socket/session"
   "log"
)

type Handler interface {
   IsValid(l *log.Logger, s *session.Session) bool

   HandleRequest(l *log.Logger, s *session.Session, r *request.RequestReader)
}
