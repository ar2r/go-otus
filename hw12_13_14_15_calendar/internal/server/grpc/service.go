package grpcserver

import (
	"context"

	dto2 "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/dto"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	pb "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc/protobuf"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/services"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

type Service struct {
	pb.UnimplementedEventServiceServer
	service services.EventServiceInterface
}

func NewService(service services.EventServiceInterface) pb.EventServiceServer {
	return &Service{
		service: service,
	}
}

func (s *Service) Create(ctx context.Context, event *pb.Event) (*pb.EventDataResponse, error) {
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
	add, err := s.service.Add(ctx, dto)
	if err != nil {
		return nil, err
	}
	event.Id = add.ID.String()
	return &pb.EventDataResponse{Event: event}, nil
}

func (s *Service) Update(ctx context.Context, event *pb.Event) (*pb.EventDataResponse, error) {
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
	add, err := s.service.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	event.Id = add.ID.String()
	return &pb.EventDataResponse{Event: event}, nil
}

func (s *Service) Delete(ctx context.Context, request *pb.DeleteEventRequest) (*pb.EmptyResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	dto := dto2.DeleteEventDto{
		ID: id,
	}
	err = s.service.Delete(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Service) ListOnDate(ctx context.Context, interval *pb.ListByDateRequest) (*pb.ListResponse, error) {
	dto := dto2.ListByDateDto{
		Date: interval.GetDate().AsTime(),
	}
	list, err := s.service.ListByDate(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.listEventsToResponse(list), nil
}

func (s *Service) ListOnWeek(ctx context.Context, interval *pb.ListByDateRequest) (*pb.ListResponse, error) {
	dto := dto2.ListByDateDto{
		Date: interval.GetDate().AsTime(),
	}
	list, err := s.service.ListByWeek(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.listEventsToResponse(list), nil
}

func (s *Service) ListOnMonth(ctx context.Context, interval *pb.ListByDateRequest) (*pb.ListResponse, error) {
	dto := dto2.ListByDateDto{Date: interval.GetDate().AsTime()}
	list, err := s.service.ListByMonth(ctx, dto)
	if err != nil {
		return nil, err
	}
	return s.listEventsToResponse(list), nil
}

func (s *Service) listEventsToResponse(list []model.Event) *pb.ListResponse {
	response := make([]*pb.Event, 0, len(list))
	for _, event := range list {
		response = append(response, &pb.Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			StartDt:     &timestamp.Timestamp{Seconds: event.StartDt.Unix()},
			EndDt:       &timestamp.Timestamp{Seconds: event.EndDt.Unix()},
			Description: &event.Description,
			UserId:      event.UserID.String(),
			NotifyAt:    &duration.Duration{Seconds: int64(event.NotifyAt)},
		})
	}
	return &pb.ListResponse{
		Events: response,
	}
}
