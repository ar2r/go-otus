package httpserver

import (
	"encoding/json"
	"net/http"

	dto2 "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/calendar/dto"
	"github.com/google/uuid"
)

// createEventHandler Обработчик создания события.
func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var reqDto dto2.CreateEventDto
	if err := json.NewDecoder(r.Body).Decode(&reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEvent, err := s.app.CreateEvent(r.Context(), reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdEvent)
}

// updateEventHandler Обработчик обновления события.
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

// deleteEventHandler Обработчик удаления события.
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

// listByDateHandler Обработчик получения списка событий на дату.
func (s *Server) listByDateHandler(w http.ResponseWriter, r *http.Request) {
	var reqDto dto2.ListByDateDto
	if err := json.NewDecoder(r.Body).Decode(&reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.ListByDate(r.Context(), reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

// listByWeekHandler Обработчик получения списка событий на неделю.
func (s *Server) listByWeekHandler(w http.ResponseWriter, r *http.Request) {
	var reqDto dto2.ListByDateDto
	if err := json.NewDecoder(r.Body).Decode(&reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.ListByWeek(r.Context(), reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

// listByMonthHandler Обработчик получения списка событий за месяц.
func (s *Server) listByMonthHandler(w http.ResponseWriter, r *http.Request) {
	var reqDto dto2.ListByDateDto
	if err := json.NewDecoder(r.Body).Decode(&reqDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.ListByMonth(r.Context(), reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}
