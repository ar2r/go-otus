package integration

import (
	"context"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage/pgx"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	TestCoreSuite
	r model.EventRepository
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupSuite() {
	s.TestCoreSuite.SetupSuite()
}

func (s *RepositoryTestSuite) SetupTest() {
	s.TestCoreSuite.SetupTest()
	s.r, _ = sqlstorage.New(s.logg, s.pool)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.TestCoreSuite.TearDownTest()
}

func (s *RepositoryTestSuite) TearDownSuite() {
	s.TestCoreSuite.TearDownSuite()
}

func (s *RepositoryTestSuite) TestSave() {
	startTime := time.Now()
	item := model.Event{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Title:       "test",
		Description: "description",
		StartDt:     startTime,
		EndDt:       startTime,
	}
	savedEvent, err := s.r.Add(context.Background(), item)
	s.Require().NoError(err)
	s.Require().NotEqual(0, savedEvent.ID)

	dbItem := s.getDirectItem("test")

	s.Require().Equal(savedEvent.ID, dbItem.ID)
	s.Require().Equal(item.Title, dbItem.Title)
	s.Require().Equal(item.Description, dbItem.Description)
	s.Require().Equal(item.StartDt.Format(time.DateTime), dbItem.StartDt.Format(time.DateTime))
	s.Require().Equal(item.EndDt.Format(time.DateTime), dbItem.EndDt.Format(time.DateTime))
}

func (s *RepositoryTestSuite) TestGet() {
	startTime := time.Now()
	endDt := startTime.Add(time.Hour)
	existItem := model.Event{
		ID:          uuid.New(),
		Title:       "test",
		Description: "description",
		StartDt:     startTime,
		EndDt:       endDt,
	}
	id := s.saveDirectItem(existItem)

	item, err := s.r.Get(context.Background(), id)
	s.Require().NoError(err)
	s.Require().NotEqual(0, item.ID)

	s.Require().Equal(id, item.ID)
	s.Require().Equal(item.Title, existItem.Title)
	s.Require().Equal(item.Description, existItem.Description)
	s.Require().Equal(item.StartDt.Format(time.DateTime), existItem.StartDt.Format(time.DateTime))
	s.Require().Equal(item.EndDt.Format(time.DateTime), existItem.EndDt.Format(time.DateTime))
}

func (s *RepositoryTestSuite) TestGetNotFound() {
	_, err := s.r.Get(context.Background(), uuid.New())
	s.Require().Error(err)
	s.Require().ErrorIs(err, pgx.ErrNoRows)
}
