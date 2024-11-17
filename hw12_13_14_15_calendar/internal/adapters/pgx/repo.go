package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/db"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model/event"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pgxPool *pgxpool.Pool
}

func connect(ctx context.Context, conf config.DatabaseConf, logg *logger.Logger) (*pgxpool.Pool, error) {
	pgxPool, err := db.Connect(ctx, conf, logg)
	if err != nil {
		logg.Error(fmt.Sprintf("failed to create connection to database: %s", err))
		return nil, err
	}
	return pgxPool, nil
}

func New(ctx context.Context, conf config.DatabaseConf, logg *logger.Logger) (*Storage, error) {
	pgxPool, err := connect(ctx, conf, logg)
	if err != nil {
		return nil, err
	}

	return &Storage{
		pgxPool: pgxPool,
	}, nil
}

func (s *Storage) Close() {
	s.pgxPool.Close()
}

// CreateEvent Создать событие с проверками на возможные пересечения с другими событиями.
func (s *Storage) CreateEvent(ctx context.Context, userID, id uuid.UUID, title string) error {
	event := event.Event{
		ID:          id,
		Title:       title,
		Description: "",
		StartDt:     time.Time{},
		EndDt:       time.Time{},
		UserID:      userID,
	}

	_, err := s.Add(ctx, event)
	return err
}

// Get Вернуть событие по идентификатору.
func (s *Storage) Get(ctx context.Context, id uuid.UUID) (*event.Event, error) {
	row := s.pgxPool.QueryRow(ctx, "SELECT * FROM events WHERE id = $1", id)

	var event event.Event
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDt,
		&event.EndDt,
		&event.UserID,
		&event.Notify)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// Add Добавить событие.
func (s *Storage) Add(ctx context.Context, event event.Event) (*event.Event, error) {
	_, err := s.pgxPool.Exec(ctx,
		"INSERT INTO events (id, title, description, start_dt, end_dt, user_id, notify) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.Title, event.Description, event.StartDt, event.EndDt, event.UserID, event.Notify)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Update Обновить событие.
func (s *Storage) Update(ctx context.Context, event event.Event) (*event.Event, error) {
	_, err := s.pgxPool.Exec(ctx,
		"UPDATE events SET title = $1, description = $2, start_dt = $3, end_dt = $4, user_id = $5, notify = $6 WHERE id = $7",
		event.Title, event.Description, event.StartDt, event.EndDt, event.UserID, event.Notify, event.ID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Delete Удалить событие.
func (s *Storage) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.pgxPool.Exec(ctx, "DELETE FROM events WHERE id = $1", uuid)
	if err != nil {
		return err
	}
	return nil
}

// List Вернуть все события.
func (s *Storage) List(ctx context.Context) ([]event.Event, error) {
	rows, err := s.pgxPool.Query(ctx, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var event event.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDt,
			&event.EndDt,
			&event.UserID,
			&event.Notify)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// ListByDate Найти все события, которые происходят в указанный день.
func (s *Storage) ListByDate(_ context.Context, _ time.Time) ([]event.Event, error) {
	return nil, adapters.ErrNotImplemented
}

// ListByPeriod Найти события, которые пересекается с указанным временным промежутком.
func (s *Storage) ListByPeriod(_ context.Context, _ time.Time, _ time.Time) ([]event.Event, error) {
	return nil, adapters.ErrNotImplemented
}

// ListByWeek Найти все события, которые происходят в указанной неделе.
// Неделя начинается с понедельника.
func (s *Storage) ListByWeek(_ context.Context, _ time.Time) ([]event.Event, error) {
	return nil, adapters.ErrNotImplemented
}

// ListByMonth Найти все события, которые происходят в указанном месяце.
// Месяц начинается с первого числа.
func (s *Storage) ListByMonth(_ context.Context, _ time.Time) ([]event.Event, error) {
	return nil, adapters.ErrNotImplemented
}
