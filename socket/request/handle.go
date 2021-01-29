package request

type Handle interface {
   Handle(int, RequestReader)
}

type Supplier interface {
   Supply(op uint16) Handle
}
