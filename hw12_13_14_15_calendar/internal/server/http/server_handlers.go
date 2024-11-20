package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var e model.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.app.CreateEvent(r.Context(), e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/events/"):]
	eventID := uuid.MustParse(id)
	e, err := s.app.GetEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/events/"):]
	eventID := uuid.MustParse(id)
	if err := s.app.DeleteEvent(r.Context(), eventID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
