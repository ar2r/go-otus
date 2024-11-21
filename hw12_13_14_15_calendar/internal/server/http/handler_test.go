package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/dto"
	mock_app "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_createEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestDto := dto.CreateEventDto{
		UserID:      uuid.New(),
		Title:       "Title",
		Description: "Description",
		StartDt:     time.Now().Truncate(time.Second),
		EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
		NotifyAt:    time.Hour,
	}
	requestBody, _ := json.Marshal(requestDto)
	responseModel := requestDto.ToModel()
	responseJson, _ := json.Marshal(responseModel)

	mockApp := mock_app.NewMockIApplication(ctrl)
	server := &Server{app: mockApp}
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	mux := server.registerRoutes()
	mockApp.EXPECT().CreateEvent(gomock.Any(), requestDto).Return(responseModel, nil)
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(responseJson), strings.TrimSpace(rr.Body.String()))
}

func TestHandler_updateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestDto := dto.UpdateEventDto{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Title:       "New Title",
		Description: "New Description",
		StartDt:     time.Now().Truncate(time.Second),
		EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
		NotifyAt:    time.Hour,
	}
	requestBody, _ := json.Marshal(requestDto)
	responseModel := requestDto.ToModel()
	responseJson, _ := json.Marshal(responseModel)

	mockApp := mock_app.NewMockIApplication(ctrl)
	server := &Server{app: mockApp}
	req, _ := http.NewRequest("PUT", "/events/"+requestDto.ID.String(), bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	mux := server.registerRoutes()
	mockApp.EXPECT().UpdateEvent(gomock.Any(), requestDto).Return(responseModel, nil)
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(responseJson), strings.TrimSpace(rr.Body.String()))
}

func TestHandler_deleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestDto := dto.DeleteEventDto{
		ID: uuid.New(),
	}
	requestBody, _ := json.Marshal(requestDto)

	mockApp := mock_app.NewMockIApplication(ctrl)
	server := &Server{app: mockApp}
	req, _ := http.NewRequest("DELETE", "/events/"+requestDto.ID.String(), bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	mux := server.registerRoutes()
	mockApp.EXPECT().DeleteEvent(gomock.Any(), requestDto).Return(nil)
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, strings.TrimSpace(rr.Body.String()))
}

func TestHandler_listByDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestDto := dto.UpdateEventDto{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Title:       "New Title",
		Description: "New Description",
		StartDt:     time.Now().Truncate(time.Second),
		EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
		NotifyAt:    time.Hour,
	}
	requestBody, _ := json.Marshal(requestDto)
	responseModel := requestDto.ToModel()
	responseJson, _ := json.Marshal(responseModel)

	mockApp := mock_app.NewMockIApplication(ctrl)
	server := &Server{app: mockApp}
	req, _ := http.NewRequest("PUT", "/events/"+requestDto.ID.String(), bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	mux := server.registerRoutes()
	mockApp.EXPECT().ListByDate(gomock.Any(), requestDto).Return(responseModel, nil)
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(responseJson), strings.TrimSpace(rr.Body.String()))
}
