package request

type Handler interface {
   Handle(int, RequestReader)
}

type Supplier interface {
   Supply(op uint16) Handler
}
