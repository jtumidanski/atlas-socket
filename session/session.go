package session

import (
   "github.com/jtumidanski/atlas-socket/crypto"
   "net"
   "time"
)

type Session interface {
   SessionId() int
   ReceiveAESOFB() *crypto.AESOFB
   WriteHello()
   LastRequest() time.Time
   UpdateLastRequest()
   GetRemoteAddress() net.Addr
}

type Creator interface {
   Create(int, net.Conn) (Session, error)
}

type Destroyer interface {
   Destroy(int)
}

type Retriever interface {
   Get(int) Session
   GetAll() []Session
}

type Service interface {
   Creator
   Retriever
   Destroyer
}
