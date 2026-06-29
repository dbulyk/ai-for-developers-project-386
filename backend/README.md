# Calendar Booking API — Backend

REST API для записи на встречи по календарю. Реализовано на Go 1.25 с использованием роутера `chi`, стандартного пакета `log/slog` и in-memory хранилища.

## Стек

- **Go** 1.25
- **Роутинг:** `github.com/go-chi/chi/v5`
- **Конфигурация:** `github.com/caarlos0/env/v11`
- **Логирование:** `log/slog` (JSON в `prod`-режиме, text в dev)
- **Тесты:** `testify`, `httptest`, table-driven + race
- **Линтер:** `golangci-lint`

## Требования

- Go >= 1.25
- `make` (опционально, для команд ниже)
- `golangci-lint` — устанавливается локально через `make tools`

## Команды Makefile

| Команда | Описание |
|---|---|
| `make run` | Запустить сервер в dev-режиме на `PORT` (по умолчанию `8080`) |
| `make build` | Собрать бинарник в `bin/server` |
| `make test` | Запустить unit- и http-тесты |
| `make lint` | Запустить `golangci-lint` |
| `make fmt` | Отформатировать код `gofmt` |
| `make vet` | Запустить `go vet ./...` |
| `make check` | Последовательно: `fmt` → `vet` → `lint` → `test` |
| `make tidy` | `go mod tidy` |
| `make clean` | Удалить `bin/` |
| `make tools` | Установить `golangci-lint` в `./bin` |
| `make smoke` | Поднять сервер и прогнать curl-сценарий из `scripts/smoke.sh` |

## Быстрый старт

```bash
cd backend
make tools      # один раз
make run        # сервер на http://localhost:8080
```

Проверка работоспособности:

```bash
curl http://localhost:8080/health
```

## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|---|---|---|
| `PORT` | Порт HTTP-сервера | `8080` |
| `LOG_FORMAT` | Формат логов: `text` или `json` | `text` |
| `LOG_LEVEL` | Уровень логирования (`debug`, `info`, `warn`, `error`) | `info` |
| `CORS_ALLOWED_ORIGINS` | Разрешённые origins через запятую | `http://localhost:4010` |
| `OWNER_TIMEZONE` | Таймзона владельца календаря для генерации слотов | `Europe/Moscow` |

## Структура проекта

```
backend/
├── cmd/server/main.go          # точка входа: config → logger → store → handlers → server
├── internal/
│   ├── config/                 # env-конфиг
│   ├── logger/                 # slog setup
│   ├── clock/                  # абстракция time.Now() + MockClock
│   ├── models/                 # EventType, Booking + валидация
│   ├── store/                  # Store interface + in-memory реализация
│   ├── slots/                  # генерация доступных слотов
│   ├── middleware/             # CORS, request logging
│   ├── handlers/               # admin/public хендлеры
│   └── server/                 # сборка chi-роутера
├── scripts/smoke.sh            # curl-smoke тест
├── Makefile
└── .golangci.yml
```

## Эндпоинты

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/public/event-types` | Список типов событий |
| `GET` | `/public/event-types/{id}/slots` | Доступные слоты |
| `POST` | `/public/bookings` | Создать бронирование |
| `GET` | `/admin/event-types` | Список типов событий |
| `POST` | `/admin/event-types` | Создать тип события |
| `PUT` | `/admin/event-types/{id}` | Обновить тип события |
| `DELETE` | `/admin/event-types/{id}` | Удалить тип события |
| `GET` | `/admin/bookings` | Предстоящие бронирования |
| `DELETE` | `/admin/bookings/{id}` | Отменить бронирование |

## Архитектурные заметки

### Слои

```
cmd/server/main.go
       │
       ▼
internal/server/server.go        ← сборка роутера, middleware, хендлеров
       │
       ├── internal/handlers      ← http-адаптеры
       │       ├── store.Store
       │       ├── clock.Clock
       │       └── slots.Generate
       │
       ├── internal/middleware    ← CORS, logging, recoverer
       │
       └── internal/config        ← env-конфиг
```

### Хранилище

- Реализация — `internal/store/memory.go`.
- In-memory; все данные теряются при перезапуске сервера (ограничение MVP).
- Потокобезопасность через `sync.RWMutex` на уровне `MemoryStore`.
- Двойное бронирование одного слота исключено: `CreateBooking` атомарно проверяет занятость и возвращает `store.ErrSlotTaken`.

### Конфигурация

- Все настройки — через переменные окружения.
- Чтобы добавить новое поле: обновить `internal/config/config.go` (теги `env` / `envDefault`), затем дополнить `internal/config/config_test.go`.

### Генерация слотов

- Рабочие дни: пн–пт.
- Рабочие часы: 09:00–18:00 в таймзоне `OWNER_TIMEZONE`.
- Шаг сетки равен `durationMinutes` типа события.
- Последний возможный старт: `18:00 - duration`. Если `duration > 9h`, слотов нет.
- Окно бронирования: 14 календарных дней **включая сегодня**.
- Прошедшие слоты текущего дня исключаются автоматически.
- Занятые слоты не отдаются публичному API.

### Логирование

- Через `log/slog`.
- Чувствительные поля (например, `guestName`) никогда не логируются.

## Smoke test

```bash
cd backend
make smoke
```

Сценарий поднимает сервер и последовательно проверяет: health, создание типа события, публичный список, слоты, создание бронирования, список бронирований, удаление бронирования, удаление типа события.

## Покрытие тестами

Последняя проверка `go test -cover ./...`:

<!-- coverage -->

| Пакет | Покрытие |
|---|---|
| `internal/clock` | 100.0% |
| `internal/slots` | 100.0% |
| `internal/store` | 98.0% |
| `internal/models` | 96.9% |
| `internal/middleware` | 96.7% |
| `internal/logger` | 88.9% |
| `internal/handlers` | 83.4% |
| `internal/server` | 81.5% |
| `internal/config` | 75.0% |
| `cmd/server` | 31.0% |

## Известные ограничения

1. **In-memory хранилище** — данные не сохраняются между перезапусками.
2. **Конкурентность** — `sync.RWMutex` на уровне store; для 7 эндпоинтов MVP достаточно, но при росте нагрузки стоит рассмотреть `sync.Map` или персистентное хранилище.
3. **DST-переходы** — при использовании таймзон с переходом на летнее/зимнее время генерация слотов может дать артефакты. Для `Europe/Moscow` (без DST) неактуально.
4. **CORS в проде** — по умолчанию разрешён `http://localhost:4010`. В production переопределите через `CORS_ALLOWED_ORIGINS`.
