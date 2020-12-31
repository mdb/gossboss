package gossboss

import (
	"encoding/json"
	"log"
	"net/http"
)

// Server is a gossboss server.
type Server struct {
	GossServers []string
	Port        string
	Client      *Client
}

// NewServer returns a new Server.
func NewServer(port string, gossServers []string) *Server {
	return &Server{
		Port:        port,
		GossServers: gossServers,
		Client:      NewClient(),
	}
}

// Serve establishes a gossboss server.
func (s *Server) Serve() {
	http.HandleFunc("/healthzs", s.HandleHealthzs)

	log.Println("Starting server on", s.Port)
	if err := http.ListenAndServe(s.Port, nil); err != nil {
		log.Fatal(err)
	}
}

// HandleHealthzs collects the /healthz responses from all the GossServers
// and returns a JSON array of their responses.
func (s *Server) HandleHealthzs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hzs := s.Client.CollectHealthzs(s.GossServers)

	if hzs.Summary.Failed != 0 || hzs.Summary.Errored != 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err := json.NewEncoder(w).Encode(hzs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
