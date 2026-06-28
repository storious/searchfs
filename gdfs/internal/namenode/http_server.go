package namenode

import (
	"encoding/json"
	"net/http"
	"strings"
)

type HTTPServer struct {
	node *NameNode
	mux  *http.ServeMux
}

func NewHTTPServer(node *NameNode) *HTTPServer {
	s := &HTTPServer{
		node: node,
		mux:  http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *HTTPServer) routes() {
	s.mux.HandleFunc("/files/", s.handleFile)
}

func (s *HTTPServer) handleFile(w http.ResponseWriter, r *http.Request) {
	path := FilePath("/" + strings.TrimPrefix(r.URL.Path, "/files/"))
	if path == "/" {
		http.Error(w, "missing file path", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		s.handlePutFile(w, r, path)
	case http.MethodGet:
		s.handleGetFile(w, r, path)
	case http.MethodDelete:
		s.handleDeleteFile(w, r, path)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *HTTPServer) handlePutFile(w http.ResponseWriter, r *http.Request, path FilePath) {
	var meta FileMetadata
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	meta.Path = path

	if err := s.node.CreateFile(r.Context(), meta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, meta)
}

func (s *HTTPServer) handleGetFile(w http.ResponseWriter, r *http.Request, path FilePath) {
	meta, err := s.node.GetFile(r.Context(), path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, meta)
}

func (s *HTTPServer) handleDeleteFile(w http.ResponseWriter, r *http.Request, path FilePath) {
	if err := s.node.DeleteFile(r.Context(), path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
