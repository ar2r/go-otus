package httpserver

import (
	"encoding/json"
	"net/http"

	dto2 "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/dto"
	"github.com/google/uuid"
)

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var reqDto dto2.CreateEventDto
	if err := json.NewDecoder(r.Body).Decode(&reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.app.CreateEvent(r.Context(), reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var reqDto dto2.UpdateEventDto
	if err := json.NewDecoder(r.Body).Decode(&reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedEvent, err := s.app.UpdateEvent(r.Context(), reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedEvent)
}
func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/events/"):]
	eventID := uuid.MustParse(id)
	reqDto := dto2.DeleteEventDto{ID: eventID}
	if err := s.app.DeleteEvent(r.Context(), reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
