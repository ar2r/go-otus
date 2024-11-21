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

	mockApp := mock_app.NewMockIApplication(ctrl)
	server := &Server{app: mockApp}

	reqDto := dto.CreateEventDto{
		UserID:      uuid.New(),
		Title:       "Title",
		Description: "Description",
		StartDt:     time.Now().Truncate(time.Second),
		EndDt:       time.Now().Truncate(time.Second).Add(time.Hour),
		NotifyAt:    time.Hour,
	}

	reqBody, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(reqBody))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.createEventHandler)

	model := reqDto.ToModel()
	expectedJson, _ := json.Marshal(model)

	mockApp.EXPECT().CreateEvent(gomock.Any(), reqDto).Return(model, nil)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(expectedJson), strings.TrimSpace(rr.Body.String()))
}

//func TestUpdateEventHandler(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockApp := NewMockIApplication(ctrl)
//	server := &Server{app: mockApp}
//
//	reqDto := dto.UpdateEventDto{
//		ID:    uuid.New(),
//		Title: "Updated Event",
//	}
//	reqBody, _ := json.Marshal(reqDto)
//	req, err := http.NewRequest("PUT", "/events/"+reqDto.ID.String(), bytes.NewBuffer(reqBody))
//	assert.NoError(t, err)
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(server.updateEventHandler)
//
//	mockApp.EXPECT().UpdateEvent(gomock.Any(), reqDto).Return(&dto.EventDto{ID: reqDto.ID}, nil)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//}
//
//func TestDeleteEventHandler(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockApp := NewMockIApplication(ctrl)
//	server := &Server{app: mockApp}
//
//	eventID := uuid.New()
//	req, err := http.NewRequest("DELETE", "/events/"+eventID.String(), nil)
//	assert.NoError(t, err)
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(server.deleteEventHandler)
//
//	mockApp.EXPECT().DeleteEvent(gomock.Any(), dto.DeleteEventDto{ID: eventID}).Return(nil)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusNoContent, rr.Code)
//}
