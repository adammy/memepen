package server

// Server defines the interface for a meme Server.
type Server interface {
	// Start the Server.
	Start() error
}
