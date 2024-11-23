package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event Событие - основная сущность, содержит в себе поля.
type Event struct {
	ID          uuid.UUID     // ID - уникальный идентификатор события (можно воспользоваться UUID);
	Title       string        // Заголовок - короткий текст;
	StartDt     time.Time     // Дата и время события;
	EndDt       time.Time     // Длительность события (или дата и время окончания);
	Description string        // Описание события - длинный текст, опционально;
	UserID      uuid.UUID     // ID пользователя, владельца события;
	NotifyAt    time.Duration // За сколько времени высылать уведомление, опционально.
}

func (e *Event) GenerateID() error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate uuid failed: %w", err)
	}
	e.ID = id
	return nil
}
