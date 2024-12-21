package integration

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	grpcprotobuf "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc/protobuf"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/myslog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApiTestSuite struct {
	suite.Suite
	pool        *pgxpool.Pool
	logg        *slog.Logger
	m           *migrate.Migrate
	eventClient grpcprotobuf.EventServiceClient
	ctx         context.Context
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) SetupSuite() {
	dsn := "://calendar:calendar-pwd@localhost:5432/calendar"
	pool, err := pgxpool.New(context.Background(), "postgres"+dsn)
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool
	s.logg = myslog.New("debug", "stdout", "")

	s.m, err = migrate.New("file://../../migrations", "pgx"+dsn)
	s.Require().NoError(err)

	// grpc
	host := "localhost:9999"
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.ctx = context.Background()
	s.eventClient = grpcprotobuf.NewEventServiceClient(conn)
}

func (s *ApiTestSuite) SetupTest() {
	err := s.m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}
}

func (s *ApiTestSuite) TearDownTest() {
	err := s.m.Down()
	s.Require().NoError(err)
}

func (s *ApiTestSuite) TearDownSuite() {
	s.m.Close()
	s.pool.Close()
}

func (s *ApiTestSuite) TestCreate() {
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

func (s *ApiTestSuite) TestListByDate() {
	// Arrange
	s.createDirectEventStartsAt(time.Now().Add(-time.Hour*24), "event за предыдущий день")
	expectedEventId := s.createDirectEventStartsAt(time.Now(), "event за текущий день")
	s.createDirectEventStartsAt(time.Now().Add(+time.Hour*24), "event за следующий день")

	// Act
	listResp, err := s.eventClient.ListByDate(s.ctx, &grpcprotobuf.ListByDateRequest{Date: timestamppb.Now()})

	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
	s.Require().Equal(expectedEventId.String(), listResp.Events[0].Id)
	s.Require().Equal(1, len(listResp.Events))
}

func (s *ApiTestSuite) TestListByWeek() {
	// Arrange
	s.createDirectEventStartsAt(time.Now().Add(-time.Hour*24*7), "event за предыдущую неделю")
	expectedEventId := s.createDirectEventStartsAt(time.Now(), "event за текущую неделю")
	s.createDirectEventStartsAt(time.Now().Add(+time.Hour*24*7), "event за следующую неделю")

	// Act
	listResp, err := s.eventClient.ListByWeek(s.ctx, &grpcprotobuf.ListByDateRequest{Date: timestamppb.Now()})

	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
	s.Require().Equal(expectedEventId.String(), listResp.Events[0].Id)
	s.Require().Equal(1, len(listResp.Events))

}

func (s *ApiTestSuite) TestListByMonth() {
	// Arrange
	s.createDirectEventStartsAt(time.Now().Add(-time.Hour*24*31), "event за предыдущий месяц")
	expectedEventId := s.createDirectEventStartsAt(time.Now(), "event за текущий месяц")
	s.createDirectEventStartsAt(time.Now().Add(+time.Hour*24*31), "event за следующий месяц")

	// Act
	listResp, err := s.eventClient.ListByMonth(s.ctx, &grpcprotobuf.ListByDateRequest{Date: timestamppb.Now()})

	// Assert
	s.Require().NoError(err)
	s.Require().NotNil(listResp)
	s.Require().NotEmpty(listResp.Events)
	s.Require().Equal(expectedEventId.String(), listResp.Events[0].Id)
	s.Require().Equal(1, len(listResp.Events))
}

func (s *ApiTestSuite) createDirectEventStartsAt(start time.Time, title string) uuid.UUID {
	item := model.Event{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   title,
		StartDt: start.Add(-time.Minute),
		EndDt:   start.Add(time.Minute),
	}
	return s.saveDirectItem(item)
}

func (s *ApiTestSuite) saveDirectItem(item model.Event) uuid.UUID {
	query, args, err := sq.
		Insert("events").
		Columns(
			"id",
			"user_id",
			"title",
			"description",
			"start_dt",
			"end_dt",
			"notify_at",
		).
		Values(
			uuid.New(),
			uuid.New(),
			item.Title,
			item.Description,
			item.StartDt,
			item.EndDt,
			time.Duration(0),
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		s.Fail(err.Error())
	}

	rows, err := s.pool.Query(context.Background(), query, args...)
	if err != nil {
		s.Fail(err.Error())
	}
	defer rows.Close()

	var itemID uuid.UUID
	for rows.Next() {
		if scanErr := rows.Scan(&itemID); scanErr != nil {
			s.Fail(scanErr.Error())
		}
	}

	return itemID
}
