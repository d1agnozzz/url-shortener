package api

import (
	"encoding/json"
	"fmt"
	"github.com/d1agnozzz/url-shortener/internal/aliaser"
	"github.com/d1agnozzz/url-shortener/internal/storage"
	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type APIServer struct {
	listenAddr string
	config     APIConfig
}

type APIConfig struct {
	Storage       storage.Storage
	Aliaser       aliaser.Aliaser
	UrlSanitizer  urlsanitizer.URLSanitizer
	MaxCollisions int
}

func NewApiServer(listenAddr string, config APIConfig) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		config:     config,
	}
}

func (s *APIServer) UrlHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		s.getURLHandler(w, req)
	case "POST":
		s.postURLHandler(w, req)
	default:
		w.Header().Add("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func (s *APIServer) getURLHandler(w http.ResponseWriter, req *http.Request) {
	alias := req.URL.Path[1:]

	if len(alias) != aliaser.AliasLen {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	urlMapping, err := s.config.Storage.GetByAlias(req.Context(), alias)

	switch err {
	case nil:
		WriteJSON(w, http.StatusOK, types.UrlResponse{
			Url: urlMapping.Url,
		})
	case pgx.ErrNoRows:
		http.NotFound(w, req)
	default:
		WriteJSON(w, http.StatusInternalServerError, types.APIError{
			Error: err.Error(),
		})

	}
}

func (s *APIServer) postURLHandler(w http.ResponseWriter, req *http.Request) {

	var r types.PostURLRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		WriteJSON(w, http.StatusBadRequest, types.APIError{
			Error: err.Error(),
		})
		return
	}

	sanitized, sanErr := s.config.UrlSanitizer.Sanitize(r.Url)

	if sanErr != nil {
		WriteJSON(w, http.StatusBadRequest, types.APIError{
			Error: sanErr.Error(),
		})
		return
	}

	// attempts to resolve collisions, if any
	for attempt := 0; attempt <= s.config.MaxCollisions; attempt++ {
		salted := fmt.Sprintf("%s::%d", sanitized.String(), attempt)
		alias := s.config.Aliaser.GenerateByStr(salted)

		toInsert := types.URLMapping{
			Url:   sanitized.String(),
			Alias: alias.String(),
		}
		retrieved, getErr := s.config.Storage.GetByAlias(req.Context(), alias.String())

		switch getErr {
		case nil:
			if retrieved.Url == sanitized.String() {
				WriteJSON(w, http.StatusOK, types.AliasResponse{
					Alias: retrieved.Alias,
				})
				return
			}

		case pgx.ErrNoRows:
			s.config.Storage.InsertURLMapping(req.Context(), toInsert)
			WriteJSON(w, http.StatusOK, types.AliasResponse{
				Alias: toInsert.Alias,
			})
			return
		default:
			WriteJSON(w, http.StatusInternalServerError, types.APIError{
				Error: getErr.Error(),
			})
			return

		}

	}

	WriteJSON(w, http.StatusInternalServerError, types.APIError{
		Error: "too much collisions",
	})

}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
