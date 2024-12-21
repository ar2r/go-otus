package integration

import (
	"context"
	"testing"
	"time"

	grpcprotobuf "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc/protobuf"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCTestSuite struct {
	TestCoreSuite
	eventClient grpcprotobuf.EventServiceClient
	ctx         context.Context
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}

func (s *GRPCTestSuite) SetupSuite() {
	s.TestCoreSuite.SetupSuite()

	// grpc
	host := "localhost:9999"
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.ctx = context.Background()
	s.eventClient = grpcprotobuf.NewEventServiceClient(conn)
}

func (s *GRPCTestSuite) SetupTest() {
	s.TestCoreSuite.SetupTest()
}

func (s *GRPCTestSuite) TearDownTest() {
	s.TestCoreSuite.TearDownTest()
}

func (s *GRPCTestSuite) TearDownSuite() {
	s.TestCoreSuite.TearDownSuite()
}

func (s *GRPCTestSuite) TestCreate() {
	createReq := &grpcprotobuf.CreateEventRequest{
		UserId:   uuid.New().String(),
		Title:    "Test Event",
		StartDt:  timestamppb.Now(),
		EndDt:    timestamppb.Now(),
		NotifyAt: &durationpb.Duration{Seconds: 3600},
	}
	createdEven, err := s.eventClient.Create(s.ctx, createReq)
	s.Require().NoError(err)
	s.Require().Equal(createReq.Title, createdEven.Title)
}

func (s *GRPCTestSuite) TestListByDate() {
	// Arrange
	s.createDirectEventStartsAt(time.Now().Add(-time.Hour*24), "event за предыдущий день")
	expectedEventID := s.createDirectEventStartsAt(time.Now(), "event за текущий день")
	s.createDirectEventStartsAt(time.Now().Add(+time.Hour*24), "event за следующий день")

	// Act
	listResp, err := s.eventClient.ListByDate(s.ctx, &grpcprotobuf.ListByDateRequest{Date: timestamppb.Now()})

	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
	s.Require().Equal(expectedEventID.String(), listResp.Events[0].Id)
	s.Require().Equal(1, len(listResp.Events))
}

func (s *GRPCTestSuite) TestListByWeek() {
	// Arrange
	s.createDirectEventStartsAt(time.Now().Add(-time.Hour*24*7), "event за предыдущую неделю")
	expectedEventID := s.createDirectEventStartsAt(time.Now(), "event за текущую неделю")
	s.createDirectEventStartsAt(time.Now().Add(+time.Hour*24*7), "event за следующую неделю")

	// Act
	listResp, err := s.eventClient.ListByWeek(s.ctx, &grpcprotobuf.ListByDateRequest{Date: timestamppb.Now()})

	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
	s.Require().Equal(expectedEventID.String(), listResp.Events[0].Id)
	s.Require().Equal(1, len(listResp.Events))
}

func (s *GRPCTestSuite) TestListByMonth() {
	// Arrange
	s.createDirectEventStartsAt(time.Now().Add(-time.Hour*24*31), "event за предыдущий месяц")
	expectedEventID := s.createDirectEventStartsAt(time.Now(), "event за текущий месяц")
	s.createDirectEventStartsAt(time.Now().Add(+time.Hour*24*31), "event за следующий месяц")

	// Act
	listResp, err := s.eventClient.ListByMonth(s.ctx, &grpcprotobuf.ListByDateRequest{Date: timestamppb.Now()})

	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
	s.Require().Equal(expectedEventID.String(), listResp.Events[0].Id)
	s.Require().Equal(1, len(listResp.Events))
}
