package socket

type ServerOpt func(s *Server)

func IpAddress(ipAddress string) func(*Server) {
	return func(s *Server) {
		s.ipAddress = ipAddress
	}
}

func Port(port int) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}
