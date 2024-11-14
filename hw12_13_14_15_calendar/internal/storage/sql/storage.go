package sqlstorage

import (
	"context"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	PgxPool *pgxpool.Pool
}

func New(pgxPool *pgxpool.Pool) *Storage {
	return &Storage{
		PgxPool: pgxPool,
	}
}

// CreateEvent Создать событие с проверками на возможные пересечения с другими событиями.
func (s *Storage) CreateEvent(ctx context.Context, id uuid.UUID, title string) error {
	event := storage.Event{
		ID:          id,
		Title:       title,
		Description: "",
		StartDt:     time.Time{},
		EndDt:       time.Time{},
		UserID:      uuid.Nil,
	}

	_, err := s.Add(ctx, event)
	return err
}

// Get Вернуть событие по идентификатору.
func (s *Storage) Get(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	row := s.PgxPool.QueryRow(ctx, "SELECT * FROM events WHERE id = $1", id)

	var event storage.Event
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
func (s *Storage) Add(ctx context.Context, event storage.Event) (*storage.Event, error) {
	_, err := s.PgxPool.Exec(ctx,
		"INSERT INTO events (id, title, description, start_dt, end_dt, user_id, notify) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.Title, event.Description, event.StartDt, event.EndDt, event.UserID, event.Notify)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Update Обновить событие.
func (s *Storage) Update(ctx context.Context, event storage.Event) (*storage.Event, error) {
	_, err := s.PgxPool.Exec(ctx,
		"UPDATE events SET title = $1, description = $2, start_dt = $3, end_dt = $4, user_id = $5, notify = $6 WHERE id = $7",
		event.Title, event.Description, event.StartDt, event.EndDt, event.UserID, event.Notify, event.ID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Delete Удалить событие.
func (s *Storage) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.PgxPool.Exec(ctx, "DELETE FROM events WHERE id = $1", uuid)
	if err != nil {
		return err
	}
	return nil
}

// List Вернуть все события.
func (s *Storage) List(ctx context.Context) ([]storage.Event, error) {
	rows, err := s.PgxPool.Query(ctx, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
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
func (s *Storage) ListByDate(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, storage.ErrNotImplemented
}

// ListByPeriod Найти события, которые пересекается с указанным временным промежутком.
func (s *Storage) ListByPeriod(_ context.Context, _ time.Time, _ time.Time) ([]storage.Event, error) {
	return nil, storage.ErrNotImplemented
}

// ListByWeek Найти все события, которые происходят в указанной неделе.
// Неделя начинается с понедельника.
func (s *Storage) ListByWeek(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, storage.ErrNotImplemented
}

// ListByMonth Найти все события, которые происходят в указанном месяце.
// Месяц начинается с первого числа.
func (s *Storage) ListByMonth(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, storage.ErrNotImplemented
}
