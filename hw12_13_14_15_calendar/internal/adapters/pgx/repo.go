package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/database"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/easylog"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	logg app.Logger
	conn *pgxpool.Pool
}

func New(ctx context.Context, conf database.Config, logg *easylog.Logger) (*Storage, error) {
	pgxPool, err := connect(ctx, conf, logg)
	if err != nil {
		return nil, err
	}

	return &Storage{
		logg: logg,
		conn: pgxPool,
	}, nil
}

func connect(ctx context.Context, conf database.Config, logg *easylog.Logger) (*pgxpool.Pool, error) {
	pgxPool, err := database.Connect(ctx, conf)
	if err != nil {
		logg.Error(fmt.Sprintf("failed to create connection to database: %s", err))
		return nil, err
	}
	return pgxPool, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}

// Add Добавить событие.
func (s *Storage) Add(ctx context.Context, event model.Event) (*model.Event, error) {
	_, err := s.conn.Exec(ctx,
		"INSERT INTO events (id, title, description, start_dt, end_dt, user_id, notify) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.Title, event.Description, event.StartDt, event.EndDt, event.UserID, event.NotifyAt)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Get Вернуть событие по идентификатору.
func (s *Storage) Get(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	row := s.conn.QueryRow(ctx, "SELECT * FROM events WHERE id = $1", id)

	var event model.Event
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDt,
		&event.EndDt,
		&event.UserID,
		&event.NotifyAt)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// Update Обновить событие.
func (s *Storage) Update(ctx context.Context, event model.Event) (*model.Event, error) {
	_, err := s.conn.Exec(ctx,
		"UPDATE events SET title = $1, description = $2, start_dt = $3, end_dt = $4, user_id = $5, notify = $6 WHERE id = $7",
		event.Title, event.Description, event.StartDt, event.EndDt, event.UserID, event.NotifyAt, event.ID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Delete Удалить событие.
func (s *Storage) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM events WHERE id = $1", uuid)
	if err != nil {
		return err
	}
	return nil
}

// List Вернуть все события.
func (s *Storage) List(ctx context.Context) ([]model.Event, error) {
	rows, err := s.conn.Query(ctx, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDt,
			&event.EndDt,
			&event.UserID,
			&event.NotifyAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// ListByPeriod Найти события, которые пересекается с указанным временным промежутком.
func (s *Storage) ListByPeriod(ctx context.Context, startDt, endDt time.Time) ([]model.Event, error) {
	rows, err := s.conn.Query(ctx, "SELECT * FROM events WHERE start_dt < $1 AND end_dt > $2", endDt, startDt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDt,
			&event.EndDt,
			&event.UserID,
			&event.NotifyAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// ListByDate Найти все события, которые происходят в указанный день.
func (s *Storage) ListByDate(ctx context.Context, date time.Time) ([]model.Event, error) {
	startDt := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endDt := startDt.AddDate(0, 0, 1)
	s.logg.Debug(fmt.Sprintf("ListByDate: startDt=%s, endDt=%s", startDt, endDt))
	return s.ListByPeriod(ctx, startDt, endDt)
}

// ListByWeek Найти все события, которые происходят в указанной неделе.
// Неделя начинается с понедельника.
func (s *Storage) ListByWeek(ctx context.Context, startDt time.Time) ([]model.Event, error) {
	startOfWeek := startDt.AddDate(0, 0, -int(startDt.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	s.logg.Debug(fmt.Sprintf("ListByWeek: startOfWeek=%s, endOfWeek=%s", startOfWeek, endOfWeek))
	return s.ListByPeriod(ctx, startOfWeek, endOfWeek)
}

// ListByMonth Найти все события, которые происходят в указанном месяце.
// Месяц начинается с первого числа.
func (s *Storage) ListByMonth(ctx context.Context, startDt time.Time) ([]model.Event, error) {
	startOfMonth := time.Date(startDt.Year(), startDt.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	s.logg.Debug(fmt.Sprintf("ListByMonth: startOfMonth=%s, endOfMonth=%s", startOfMonth, endOfMonth))
	return s.ListByPeriod(ctx, startOfMonth, endOfMonth)
}
