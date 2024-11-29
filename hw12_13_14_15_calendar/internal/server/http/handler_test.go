package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/calendar/dto"
	mock_app "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/calendar/mocks"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	mockApp := mock_app.NewIApplication(t)
	server := &Server{app: mockApp}
	mux := server.registerRoutes()

	// Create
	createEvenDto := dto.CreateEventDto{
		UserID:      uuid.New(),
		Title:       "Title",
		Description: "Description",
		StartDt:     time.Now().Truncate(time.Second),
		EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
		NotifyAt:    time.Hour,
	}
	createdEvent := createEvenDto.ToModel()

	// Update
	updateEventDto := dto.UpdateEventDto{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Title:       "New Title",
		Description: "New Description",
		StartDt:     time.Now().Truncate(time.Second),
		EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
		NotifyAt:    time.Hour,
	}
	updatedEvent := updateEventDto.ToModel()

	// Delete
	deleteEventDto := dto.DeleteEventDto{
		ID: uuid.New(),
	}

	tests := []struct {
		name         string
		method       string
		url          string
		requestDto   interface{}
		mockCall     func()
		expectedCode int
		expectedBody string
	}{
		{
			name:       "create event",
			method:     "POST",
			url:        "/events",
			requestDto: createEvenDto,
			mockCall: func() {
				mockApp.
					On("CreateEvent", mock.Anything, createEvenDto).
					Once().
					Return(createdEvent, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() string {
				responseJSON, _ := json.Marshal(createdEvent)
				return string(responseJSON)
			}(),
		},
		{
			name:       "update event",
			method:     "PUT",
			url:        "/events/" + updateEventDto.ID.String(),
			requestDto: updateEventDto,
			mockCall: func() {
				mockApp.
					On("UpdateEvent", mock.Anything, updateEventDto).
					Once().
					Return(updatedEvent, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func() string {
				responseJSON, _ := json.Marshal(updatedEvent)
				return string(responseJSON)
			}(),
		},
		{
			name:       "delete event",
			method:     "DELETE",
			url:        "/events/" + deleteEventDto.ID.String(),
			requestDto: deleteEventDto,
			mockCall: func() {
				mockApp.
					On("DeleteEvent", mock.Anything, deleteEventDto).
					Once().
					Return(nil)
			},
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.requestDto)
			req, _ := http.NewRequestWithContext(context.Background(), tt.method, tt.url, bytes.NewBuffer(requestBody))
			rr := httptest.NewRecorder()
			tt.mockCall()
			mux.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestHandler_listEvents(t *testing.T) {
	requestDto := dto.ListByDateDto{Date: time.Now().Truncate(time.Second)}
	requestBody, _ := json.Marshal(requestDto)

	events := []model.Event{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Title:       "Title",
			Description: "Description",
			StartDt:     time.Now().Truncate(time.Second),
			EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
			NotifyAt:    time.Hour,
		},
	}

	responseJSON, _ := json.Marshal(events)
	mockApp := mock_app.NewIApplication(t)
	server := &Server{app: mockApp}
	mux := server.registerRoutes()

	tests := []struct {
		name           string
		url            string
		mockMethod     func()
		expectedStatus int
	}{
		{
			name: "listByDate",
			url:  "/events/day",
			mockMethod: func() {
				mockApp.
					On("ListByDate", mock.Anything, requestDto).
					Once().
					Return(events, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "listByWeek",
			url:  "/events/week",
			mockMethod: func() {
				mockApp.
					On("ListByWeek", mock.Anything, requestDto).
					Once().
					Return(events, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "listByMonth",
			url:  "/events/month",
			mockMethod: func() {
				mockApp.
					On("ListByMonth", mock.Anything, requestDto).
					Once().
					Return(events, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequestWithContext(context.Background(), "GET", tt.url, bytes.NewBuffer(requestBody))
			rr := httptest.NewRecorder()
			tt.mockMethod()
			mux.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, string(responseJSON), strings.TrimSpace(rr.Body.String()))
		})
	}
}
