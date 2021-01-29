package session

import (
	"github.com/jtumidanski/atlas-socket/crypto"
	"log"
	"net"
)

type Session interface {
	SessionId() int
	ReceiveAESOFB() *crypto.AESOFB
	WriteHello()
	Disconnect()
}

type Creator interface {
	Create(*log.Logger, int, net.Conn) (Session, error)
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
