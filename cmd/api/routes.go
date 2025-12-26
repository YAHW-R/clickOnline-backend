package clickOnline

func (s *ApiServer) route() {
	s.router.HandleFunc("/ws", s.serveWs)
}
