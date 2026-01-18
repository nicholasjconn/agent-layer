package config

// AppliesToClient reports whether the server is enabled for the given client.
func (s MCPServer) AppliesToClient(client string) bool {
	if len(s.Clients) == 0 {
		return true
	}
	for _, c := range s.Clients {
		if c == client {
			return true
		}
	}
	return false
}
