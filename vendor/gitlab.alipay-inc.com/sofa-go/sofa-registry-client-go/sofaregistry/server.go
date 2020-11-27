package sofaregistry

type Server struct {
	address string
}

func (s *Server) GetAddress() string {
	return s.address
}
