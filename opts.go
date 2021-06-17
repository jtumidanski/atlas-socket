package socket

type ServerConfigurator func(s *serverConfiguration)

func SetIpAddress(ipAddress string) func(*serverConfiguration) {
	return func(s *serverConfiguration) {
		s.ipAddress = ipAddress
	}
}

func SetPort(port int) func(*serverConfiguration) {
	return func(s *serverConfiguration) {
		s.port = port
	}
}

func SetSessionCreator(creator SessionCreator) ServerConfigurator {
	return func(s *serverConfiguration) {
		s.creator = creator
	}
}

func SetSessionDestroyer(destroyer SessionDestroyer) ServerConfigurator {
	return func(s *serverConfiguration) {
		s.destroyer = destroyer
	}
}

func SetSessionMessageDecryptor(decryptor SessionMessageDecryptor) ServerConfigurator {
	return func(s *serverConfiguration) {
		s.decryptor = decryptor
	}
}
