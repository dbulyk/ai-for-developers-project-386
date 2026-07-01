### Hexlet tests and linter status:
[![Actions Status](https://github.com/dbulyk/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/dbulyk/ai-for-developers-project-386/actions)
[![Latest Release](https://img.shields.io/github/v/release/dbulyk/ai-for-developers-project-386)](https://github.com/dbulyk/ai-for-developers-project-386/releases/latest)

# Calendar Booking API

Сервис записи на встречи по календарю. Включает:

- **REST API** для управления типами событий и бронированиями слотов.
- **Админ-панель** для создания и редактирования типов событий и просмотра записей.
- **Публичную страницу** для выбора свободного слота и записи на встречу.

Проект реализован как монорепозиторий с чётким разделением API-контракта, бэкенда и фронтенда.

## Стек

- **Бэкенд:** Go 1.25, `chi` (роутинг), `slog` (логирование), `caarlos0/env/v11` (конфиг), in-memory хранилище с `sync.RWMutex`.
- **Фронтенд:** React 19, TypeScript, Vite 8, Mantine 9, TanStack Query, `openapi-fetch`.
- **API-контракт:** TypeSpec → OpenAPI 3.
- **Тесты:** `testify` (backend), Vitest + Testing Library (frontend).
- **Линтеры:** `golangci-lint` (backend), `oxlint` (frontend).

## Структура репозитория

```
.
├── backend/   # Go-сервер
├── frontend/  # React-приложение
├── typespec/  # Спецификация API на TypeSpec
├── plans/     # Планы реализации
└── .github/   # GitHub Actions workflows
```

## Требования

- Go >= 1.25
- Node.js >= 22 (используется в CI; для локальной разработки подойдёт >= 20)
- npm

## Быстрый старт

### Вариант 1: с реальным бэкендом

1. Сгенерируйте OpenAPI-спецификацию:

   ```bash
   cd typespec
   npm install
   npm run compile
   ```

2. Запустите бэкенд на `http://localhost:8080`:

   ```bash
   cd ../backend
   make tools    # один раз: устанавливает golangci-lint в ./bin
   make run
   ```

3. Запустите фронтенд на `http://localhost:5173`:

   ```bash
   cd ../frontend
   cp .env.example .env
   npm install
   npm run dev
   ```

> **Важно:** по умолчанию бэкенд разрешает CORS только для `http://localhost:4010`. Для работы с фронтендом на `http://localhost:5173` запускайте бэкенд с переменной:
> ```bash
> CORS_ALLOWED_ORIGINS=http://localhost:5173 make run
> ```

### Вариант 2: с моком API (Prism)

Если бэкенд не нужен, можно запустить только фронтенд с фейковым API:

```bash
cd typespec
npm install
npm run compile

cd ../frontend
cp .env.example .env
npm install
npm run dev:mock
```

- Vite откроет приложение на `http://localhost:5173`.
- Prism поднимет динамический мок API на `http://localhost:4010`.

## Скрипты

### Backend

| Команда | Описание |
|---------|----------|
| `make run` | Запустить сервер |
| `make build` | Собрать бинарник в `bin/server` |
| `make test` | Запустить тесты |
| `make lint` | Запустить `golangci-lint` |
| `make fmt` | Отформатировать код (`gofmt`) |
| `make vet` | Запустить `go vet` |
| `make check` | `fmt` + `vet` + `lint` + `test` |
| `make tools` | Установить `golangci-lint` в `./bin` |
| `make tidy` | Очистить `go.mod` и `go.sum` |

### Frontend

| Команда | Описание |
|---------|----------|
| `npm run dev` | Vite dev-сервер |
| `npm run build` | Production-сборка (`tsc -b && vite build`) |
| `npm run preview` | Просмотр production-сборки |
| `npm run lint` | Запустить `oxlint` |
| `npm run typecheck` | Проверить типы (`tsc -b`) |
| `npm test` | Запустить тесты Vitest |
| `npm run gen:api` | Сгенерировать TS-типы из OpenAPI в `src/api/schema.d.ts` |
| `npm run mock` | Запустить Prism-мок API на `localhost:4010` |
| `npm run dev:mock` | Vite и Prism-мок параллельно |

### TypeSpec

| Команда | Описание |
|---------|----------|
| `npm run compile` | Собрать OpenAPI-спецификацию в `tsp-output/openapi.yaml` |

## Генерация API-клиента

После изменения спеков в `typespec/`:

```bash
cd typespec && npm run compile
cd ../frontend && npm run gen:api
```

Результат — обновлённый файл `frontend/src/api/schema.d.ts`.

## Переменные окружения

### Backend

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `PORT` | `8080` | Порт сервера |
| `LOG_FORMAT` | `text` | Формат логов: `text` или `json` |
| `LOG_LEVEL` | `info` | Уровень логирования |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:4010` | Разрешённые CORS-источники |
| `OWNER_TIMEZONE` | `Europe/Moscow` | Таймзона владельца календаря |

### Frontend

Переименуйте `frontend/.env.example` в `frontend/.env` и при необходимости измените:

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `VITE_API_BASE_URL` | `http://localhost:4010` | Базовый URL API |
| `VITE_OWNER_TIMEZONE` | `Europe/Moscow` | Таймзона владельца календаря |

## Деплой

Приложение поставляется как единый Docker-образ, в котором бэкенд раздаёт собранный фронтенд.

```bash
# Собрать образ
docker build -t calendar-booking .

# Запустить на порту из переменной PORT
docker run --rm -p 8080:8080 -e PORT=8080 calendar-booking
```

Для удобства используйте `Makefile` в корне:

```bash
make build        # собрать фронтенд и бэкенд локально
make docker-build # собрать Docker-образ
make docker-run   # запустить контейнер на $PORT (по умолчанию 8080)
```

> **Важно:** в образе SPA и API работают на одном origin, поэтому CORS-настройка `CORS_ALLOWED_ORIGINS` актуальна в основном для локальной разработки.

## Тесты и линтеры

Перед коммитом рекомендуется запустить полную проверку:

```bash
# backend
cd backend && make check

# frontend
cd frontend && npm run lint && npm run typecheck && npm test

# typespec
cd typespec && npm run compile
```

## Релизы

Проект использует [release-please](https://github.com/googleapis/release-please) и [Conventional Commits](https://www.conventionalcommits.org/).

### Формат коммитов

Каждый коммит должен соответствовать Conventional Commits:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Допустимые типы:

- `feat` — новая функциональность
- `fix` — исправление ошибки
- `perf` — улучшение производительности
- `docs` — документация
- `refactor` — рефакторинг без изменения поведения
- `test` — тесты
- `build` — сборка, зависимости
- `ci` — CI/CD
- `chore` — рутинные изменения

Для breaking changes используйте `!` после типа/scope или `BREAKING CHANGE:` в footer:

```
feat(api)!: change booking response format
```

### Как выпускается релиз

1. После мёрджа PR в `main` workflow `release-please` сканирует историю коммитов и создаёт/обновляет release-PR.
2. Release-PR содержит обновлённый `CHANGELOG.md` и предлагаемую версию по SemVer.
3. После мёрджа release-PR создаётся GitHub Release и git-тег `vX.Y.Z`.

### Hotfix

Если в уже открытый release-PR нужно добавить срочное исправление, просто смержите новый PR в `main` — release-please автоматически допишет изменения в существующий release-PR.

### Полезные ссылки

- [Conventional Commits](https://www.conventionalcommits.org/)
- [История релизов](https://github.com/dbulyk/ai-for-developers-project-386/releases)

## CI

- `.github/workflows/ci.yml` — основной пайплайн: линтер и тесты бэкенда, линтер, typecheck и тесты фронтенда, компиляция TypeSpec.
- `.github/workflows/hexlet-check.yml` — проверки от Hexlet.
- `.github/workflows/release-please.yml` — автоматический релиз по Conventional Commits.
