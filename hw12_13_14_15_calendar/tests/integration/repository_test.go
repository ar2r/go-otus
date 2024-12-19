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
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage/pgx"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/myslog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	logg *slog.Logger
	m    *migrate.Migrate
	r    model.EventRepository
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupSuite() {
	dsn := "://calendar:calendar-pwd@localhost:5432/calendar"
	pool, err := pgxpool.New(context.Background(), "postgres"+dsn)
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool
	s.logg = myslog.New("debug", "stdout", "")

	s.m, err = migrate.New("file://../../migrations", "pgx"+dsn)
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) SetupTest() {
	err := s.m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}

	s.r, _ = sqlstorage.New(s.logg, s.pool)
}

func (s *RepositoryTestSuite) TearDownTest() {
	err := s.m.Down()
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	s.m.Close()
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

func (s *RepositoryTestSuite) saveDirectItem(item model.Event) uuid.UUID {
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

func (s *RepositoryTestSuite) getDirectItem(name string) model.Event {
	query, args, err := sq.
		Select("id", "title", "description", "start_dt", "end_dt").
		From("events").
		Where(sq.Eq{"title": name}).
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
	item := model.Event{}

	for rows.Next() {
		scanErr := rows.Scan(&item.ID, &item.Title, &item.Description, &item.StartDt, &item.EndDt)
		if scanErr != nil {
			s.Fail(scanErr.Error())
		}
	}

	return item
}
