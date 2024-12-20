## Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»

Данное задание состоит из двух частей.

Не забываем про https://github.com/golang-standards/project-layout.

### 1) Докеризация сервиса
Необходимо:

* 🟢 создать Dockerfile для каждого из процессов (Календарь, Рассыльщик, Планировщик);
* 🟢 собрать образы и проверить их локальный запуск;
* 🟢 создать docker-compose файл, который запускает PostgreSQL, RabbitMQ и все микросервисы вместе
(для "неродных" сервисов использовать официальные образы из Docker Hub);
* 🟠 при желании доработать конфигурацию так, чтобы она поддерживала переменные окружения
(если вы используете библиотеку, то скорее всего она уже это умеет); в противном случае
придется "подкладывать" конфиг сервису с помощью Dockerfile / docker-compose -
при этом можно "заполнять" конфигурационный файл из переменных окружения, например
```bash
$ envsubst < config_template.json > config.json
```

* 🟢 если миграции выполняются руками, а не на старте сервиса, то также в docker-compose
должен запускаться one-shot скрипт, который делает это (применяет SQL миграции,
создавая структуру БД).
* 🟢 порты серверов, предоставляющих API, пробросить на host.

🟢 У преподавателя должна быть возможность запустить весь проект с помощью команды
`make up` (внутри `docker-compose up`) и погасить с помощью `make down`.

HTTP API, например, после запуска должно быть доступно по URL **http://localhost:8888/**.

### 2) Интеграционное тестирование
Необходимо:

* 🟢 создать отдельный пакет для интеграционных тестов.
* 🟢 реализовать интеграционные тесты на языке Go; при желании можно использовать
[godog](https://github.com/cucumber/godog) / [ginkgo](https://github.com/onsi/ginkgo), но
обязательным требованием это **не является**.
* 🟢 создать docker-compose файл, поднимающий все сервисы проекта + контейнер с интеграционными тестами;
* 🟢 расширить Makefile командой `integration-tests`, `make integration-tests` будет запускать интеграционные тесты;
**не стоит смешивать это с `make test`, иначе CI-пайплайн не пройдёт.**
* 🟢 прикрепить в Merge Request вывод команды `make integration-tests`.

Преподаватель может запустить интеграционные тесты с помощью команды `make integration-tests`:

- 🟢 команда должна поднять окружение (`docker-compose`), прогнать тесты и подчистить окружение за собой;
- 🟢 в случае успешного выполнения команда должна возвращать 0, иначе 1.

### Критерии оценки

- 🟢 Проект полностью запускается и останавливается с помощью `make up` / `make down` - 3 балла
- 🟢 Интеграционные тесты запускаются с помощью `make integration-tests`. Команда возвращает верный код ответа - 1 балл
- 🟢 Интеграционные тесты покрывают бизнес сценарии:
  - 🟢 добавление события и обработка бизнес ошибок - 2 балла
  - 🟢 получение листинга событий на день/неделю/месяц - 2 балла
  - 🔴 отправка уведомлений (необходимо доработать sender так, чтобы он информировал куда-то о статусе уведомления (
    БД/кролик)) - 2 балла

#### Зачёт от 7 баллов
