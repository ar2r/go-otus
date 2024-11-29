-- Добавляем колонку для признака уведомления
ALTER TABLE events
    ADD COLUMN notification_sent BOOLEAN DEFAULT FALSE;

-- Создаем индекс для быстрого поиска событий, у которых уведомление не выслано
CREATE INDEX idx_notification_not_sent ON events (notification_sent) WHERE notification_sent = FALSE;
