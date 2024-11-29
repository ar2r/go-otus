-- Удаляем индекс
DROP INDEX IF EXISTS idx_notification_not_sent;

-- Удаляем колонку
ALTER TABLE events
    DROP COLUMN IF EXISTS notification_sent;