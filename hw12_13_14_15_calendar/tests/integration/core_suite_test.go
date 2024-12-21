package integration

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/myslog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestCoreSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	logg *slog.Logger
	m    *migrate.Migrate
}

func (s *TestCoreSuite) SetupSuite() {
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

func (s *TestCoreSuite) SetupTest() {
	err := s.m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}
}

func (s *TestCoreSuite) TearDownTest() {
	err := s.m.Down()
	s.Require().NoError(err)
}

func (s *TestCoreSuite) TearDownSuite() {
	s.m.Close()
	s.pool.Close()
}

func (s *TestCoreSuite) saveDirectItem(item model.Event) uuid.UUID {
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

func (s *TestCoreSuite) createDirectEventStartsAt(start time.Time, title string) uuid.UUID {
	item := model.Event{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   title,
		StartDt: start.Add(-time.Minute),
		EndDt:   start.Add(time.Minute),
	}
	return s.saveDirectItem(item)
}

func (s *TestCoreSuite) getDirectItem(name string) model.Event {
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
