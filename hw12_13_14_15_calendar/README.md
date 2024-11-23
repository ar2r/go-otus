# Приложение "Календарь"

Основные возможности:

- Добавление событий
- Просмотр событий
- Удаление событий
- Отправка уведомлений перед наступлением события

## Информация для разработчиков

Выполните команды для запуска:

```shell
go mod download
cp config/config.toml.example config/config.toml
make migrate
docker compose up
```

- http://localhost:8080 - REST API
- http://localhost:8090 - GRPC API
- http://localhost:15672 - RabbitMQ Management (guest/guest)

#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:

- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

#### Ветки при выполнении

- `hw12_calendar` (от `master`) -> Merge Request в `master`
- `hw13_calendar` (от `hw12_calendar`) -> Merge Request в `hw12_calendar` (если уже вмержена, то в `master`)
- `hw14_calendar` (от `hw13_calendar`) -> Merge Request в `hw13_calendar` (если уже вмержена, то в `master`)
- `hw15_calendar` (от `hw14_calendar`) -> Merge Request в `hw14_calendar` (если уже вмержена, то в `master`)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**
