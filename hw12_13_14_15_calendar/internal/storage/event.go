package storage

import (
	"time"

	"github.com/google/uuid"
)

// Событие - основная сущность, содержит в себе поля:

type Event struct {
	Id          uuid.UUID     //Id - уникальный идентификатор события (можно воспользоваться UUID);
	Title       string        //Заголовок - короткий текст;
	StartDt     time.Time     //Дата и время события;
	EndDt       time.Time     //Длительность события (или дата и время окончания);
	Description string        //Описание события - длинный текст, опционально;
	UserId      uuid.UUID     //Id пользователя, владельца события;
	Notify      time.Duration //За сколько времени высылать уведомление, опционально.
}
