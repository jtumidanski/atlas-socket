package request

type Handler interface {
   Handle(int, RequestReader)
}
