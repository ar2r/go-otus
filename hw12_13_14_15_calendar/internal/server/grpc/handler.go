package grpcserver

import (
	"context"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	dto2 "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/dto"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	pb "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc/protobuf"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

// EventService Слой преобразования запроса pb в DTO и вызов сервис слоя приложения (Application).
type EventService struct {
	pb.UnimplementedEventServiceServer
	app app.IApplication
}

func NewService(app app.IApplication) pb.EventServiceServer {
	return &EventService{
		app: app,
	}
}

// Create Создание события.
func (s *EventService) Create(ctx context.Context, event *pb.CreateEventRequest) (*pb.EventResponse, error) {
	userID, err := uuid.Parse(event.GetUserId())
	if err != nil {
		return nil, err
	}
	dto := dto2.CreateEventDto{
		UserID:      userID,
		Title:       event.GetTitle(),
		StartDt:     event.GetStartDt().AsTime(),
		EndDt:       event.GetEndDt().AsTime(),
		Description: event.GetDescription(),
		NotifyAt:    event.GetNotifyAt().AsDuration(),
	}
	createdEvent, err := s.app.CreateEvent(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &pb.EventResponse{
		Id:          createdEvent.ID.String(),
		Title:       createdEvent.Title,
		StartDt:     &timestamp.Timestamp{Seconds: createdEvent.StartDt.Unix()},
		EndDt:       &timestamp.Timestamp{Seconds: createdEvent.EndDt.Unix()},
		Description: &createdEvent.Description,
		UserId:      createdEvent.UserID.String(),
		NotifyAt:    &duration.Duration{Seconds: int64(createdEvent.NotifyAt)},
	}, nil
}

// Update Обновление события.
func (s *EventService) Update(ctx context.Context, event *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	id, err := uuid.Parse(event.GetId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(event.GetUserId())
	if err != nil {
		return nil, err
	}
	dto := dto2.UpdateEventDto{
		ID:          id,
		Title:       event.GetTitle(),
		StartDt:     event.GetStartDt().AsTime(),
		EndDt:       event.GetEndDt().AsTime(),
		Description: event.GetDescription(),
		UserID:      userID,
		NotifyAt:    event.GetNotifyAt().AsDuration(),
	}
	updatedEvent, err := s.app.UpdateEvent(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &pb.EventResponse{
		Id:          updatedEvent.ID.String(),
		Title:       updatedEvent.Title,
		StartDt:     &timestamp.Timestamp{Seconds: updatedEvent.StartDt.Unix()},
		EndDt:       &timestamp.Timestamp{Seconds: updatedEvent.EndDt.Unix()},
		Description: &updatedEvent.Description,
		UserId:      updatedEvent.UserID.String(),
		NotifyAt:    &duration.Duration{Seconds: int64(updatedEvent.NotifyAt)},
	}, nil
}

// Delete Удаление события.
func (s *EventService) Delete(ctx context.Context, request *pb.DeleteEventRequest) (*pb.EmptyResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	dto := dto2.DeleteEventDto{ID: id}
	if err = s.app.DeleteEvent(ctx, dto); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

// ListByDate Получение списка событий на дату.
func (s *EventService) ListByDate(ctx context.Context, interval *pb.ListByDateRequest) (*pb.EventListResponse, error) {
	dto := dto2.ListByDateDto{
		Date: interval.GetDate().AsTime(),
	}
	list, err := s.app.ListByDate(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.listEventsToResponse(list), nil
}

// ListByWeek Получение списка событий на неделю.
func (s *EventService) ListByWeek(ctx context.Context, interval *pb.ListByDateRequest) (*pb.EventListResponse, error) {
	dto := dto2.ListByDateDto{
		Date: interval.GetDate().AsTime(),
	}
	list, err := s.app.ListByWeek(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.listEventsToResponse(list), nil
}

// ListByMonth Получение списка событий на месяц.
func (s *EventService) ListByMonth(ctx context.Context, interval *pb.ListByDateRequest) (*pb.EventListResponse, error) {
	dto := dto2.ListByDateDto{Date: interval.GetDate().AsTime()}
	list, err := s.app.ListByMonth(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.listEventsToResponse(list), nil
}

// listEventsToResponse Преобразование списка событий в ответ.
func (s *EventService) listEventsToResponse(list []model.Event) *pb.EventListResponse {
	response := make([]*pb.EventResponse, 0, len(list))
	for _, event := range list {
		response = append(response, &pb.EventResponse{
			Id:          event.ID.String(),
			Title:       event.Title,
			StartDt:     &timestamp.Timestamp{Seconds: event.StartDt.Unix()},
			EndDt:       &timestamp.Timestamp{Seconds: event.EndDt.Unix()},
			Description: &event.Description,
			UserId:      event.UserID.String(),
			NotifyAt:    &duration.Duration{Seconds: int64(event.NotifyAt)},
		})
	}
	return &pb.EventListResponse{
		Events: response,
	}
}
