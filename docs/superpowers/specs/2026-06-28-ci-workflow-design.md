# Дизайн CI workflow для тестов и линтера

## Статус

Утверждено. Реализовано.

## Контекст

Проект состоит из трёх компонентов:

- `backend/` — Go (1.25). `Makefile` предоставляет цели `lint` (`golangci-lint`), `test`, `check`, `tools`.
- `frontend/` — React + TypeScript + Vite. `package.json` содержит скрипты `lint` (oxlint), `test` (vitest run), `typecheck`.
- `typespec/` — TypeSpec. `package.json` содержит скрипт `compile`.

Существующий `.github/workflows/hexlet-check.yml` является авто-сгенерированным Hexlet-файлом и содержит предупреждение «DO NOT DELETE OR EDIT THIS FILE». Поэтому для собственных проверок создаётся отдельный workflow.

## Решение

Создать файл `.github/workflows/ci.yml` с одним job `check`, который последовательно запускает тесты и линтер для всех компонентов.

Выбран подход с одним job и последовательными шагами (вариант A) из-за простоты YAML, отсутствия дублирования setup-шагов и удобства чтения.

## Workflow

### Триггеры

- `push` на любую ветку.
- `pull_request` на любую ветку.

### Job: check

- **runner:** `ubuntu-latest`

### Шаги

1. **Checkout:** `actions/checkout@v6`
2. **Setup Go:** `actions/setup-go@v5`, версия 1.25 (из `backend/go.mod`)
3. **Setup Node:** `actions/setup-node@v4`, Node 22 (LTS), с кешем по `package-lock.json` в `frontend/` и `typespec/`
4. **Backend:**
   - `make -C backend tools` — установка `golangci-lint` в `backend/bin/`
   - `make -C backend lint`
   - `make -C backend test`
5. **Frontend:**
   - `npm ci --prefix frontend --legacy-peer-deps` — флаг нужен из-за конфликта peer-зависимостей (`openapi-typescript@7.13.0` требует `typescript@^5.x`, в проекте `typescript@~6.0.2`)
   - `npm run lint --prefix frontend`
   - `npm run typecheck --prefix frontend`
   - `npm run test --prefix frontend`
6. **Typespec:**
   - `npm ci --prefix typespec`
   - `npm run compile --prefix typespec`

## Обработка ошибок

Любой шаг, завершившийся ненулевым кодом, прерывает job и помечает workflow как failed. Это блокирует merge PR, если в репозитории настроены защитные правила.

## Исключённые альтернативы

- **Параллельные jobs по компонентам:** отклонено, так как для текущего размера проекта усложнение YAML и дублирование setup-шагов не окупается выигрышем во времени.
- **Matrix-стратегия:** отклонено, так как компоненты имеют разные требования к окружению и шагам, что делает матрицу менее читаемой.

## Примечания

- `hexlet-check.yml` не изменяется.
- Для backend требуется предварительный `make tools`, потому что `make lint` вызывает локальный бинарник `backend/bin/golangci-lint`.
