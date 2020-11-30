package sofabolt

type Handler interface {
	ServeSofaBOLT(rw ResponseWriter, req *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (s HandlerFunc) ServeSofaBOLT(rw ResponseWriter, req *Request) {
	s(rw, req)
}
