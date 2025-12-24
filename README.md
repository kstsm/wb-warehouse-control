# WarehouseControl

## Описание

Сервис предназначен для учёта товаров на складе с поддержкой истории изменений, ролевой модели доступа и аудита всех операций.

История изменений реализована через триггеры PostgreSQL для демонстрации их использования и недостатков. 

**Основные возможности:**
- CRUD-операции с товарами на складе
- История всех изменений с сохранением старых и новых значений
- Ролевая модель доступа (admin, manager, viewer)
- JWT авторизация
- Просмотр различий между версиями товаров
- Фильтрация и поиск истории изменений
- Экспорт истории в CSV
- Веб-интерфейс для управления товарами и просмотра истории

## HTTP API

- POST /api/login - авторизация и получение JWT токена
- POST /api/items - создание товара (требует роль admin или manager)
- GET /api/items - получение списка товаров
- GET /api/items/{id} - получение товара по ID
- PUT /api/items/{id} - обновление товара (требует роль admin или manager)
- DELETE /api/items/{id} - удаление товара (требует роль admin)
- GET /api/items/{id}/history - получение истории изменений товара
- GET /api/history - получение истории с фильтрами
- GET /api/history/export - экспорт истории в CSV

## Роли пользователей

- **admin** - полный доступ ко всем операциям
- **manager** - может просматривать и редактировать товары
- **viewer** - только просмотр товаров и истории

## Установка и запуск проекта

### 1. Клонирование репозитория

```bash
git clone https://github.com/kstsm/wb-warehouse-control
```

### 2. Настройка переменных окружения

Создайте `.env` файл, скопировав в него значения из `.example.env`:

```bash
cp .example.env .env
```

Отредактируйте `.env` файл, указав необходимые значения:

### 3. Запуск зависимостей (Docker)

```bash
make up
```

Это запустит PostgreSQL в контейнере Docker.

### 4. Миграция базы данных

```bash
make migrate-up
```

### 5. Запуск сервиса

```bash
make run
```
Сервис будет доступен по адресу: http://localhost:8080
___
## Линтер

Проект использует **golangci-lint** для проверки качества кода. 

### Запуск линтера

```bash
make linter
```
___
# API запросы

## POST /api/login - Авторизация

**URL:** `http://localhost:8080/api/login`

**Content-Type:** `application/json`

**Параметры:**

- `user_name` (обязательно) - имя пользователя (только буквы)
- `role` (обязательно) - роль пользователя: "admin", "manager" или "viewer"

**Body:**

```json
{
  "user_name": "admin",
  "role": "admin"
}
```

**Ожидаемый ответ (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "admin"
}
```

### Ошибки:

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Role' failed on the 'role' tag"
}
```

**Пользователь уже существует с другой ролью (409 Conflict):**

```json
{
  "error": "user already exists with role"
}
```

---

## POST /api/items - Создание товара

**URL:** `http://localhost:8080/api/items`

**Content-Type:** `application/json`

**Authorization:** `Bearer {token}` (требует роль admin или manager)

**Параметры:**

- `name` (обязательно) - название товара (минимум 1 символ)
- `description` (опционально) - описание товара
- `quantity` (обязательно) - количество товара (минимум 0)
- `price` (обязательно) - цена в копейках (минимум 0, максимум 2147483647)

**Body:**

```json
{
  "name": "Видеокарта",
  "description": "Palit GeForce RTX 5090 GameRock OC",
  "quantity": 100,
  "price": 31399900
}
```

**Ожидаемый ответ (201 Created):**

```json
{
  "id": "3f60f58c-de48-4990-9ba0-17a3f9684b4c",
  "name": "Видеокарта",
  "description": "Palit GeForce RTX 5090 GameRock OC",
  "quantity": 100,
  "price": "313999.00",
  "created_at": "2025-12-24T18:51:02Z",
  "updated_at": "2025-12-24T18:51:02Z",
  "message": "item created successfully"
}
```

### Ошибки:

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Name' failed on the 'required' tag"
}
```

```json
{
  "error": "validation for 'Quantity' failed on the 'min' tag"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Доступ запрещен (403 Forbidden):**

```json
{
  "error": "forbidden"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/items/{id} - Получение товара по ID

**URL:** `http://localhost:8080/api/items/{id}`

**Authorization:** `Bearer {token}`

**Параметры:**

- `{id}` (обязательно) - UUID товара

**Ожидаемый ответ (200 OK):**

```json
{
  "id": "7a987f51-4841-4d88-ab56-03f3c2fd7d1e",
  "name": "Видеокарта",
  "description": "Palit GeForce RTX 5090 GameRock OC",
  "quantity": 100,
  "price": "313999.00",
  "created_at": "2025-12-24T18:58:18Z",
  "updated_at": "2025-12-24T18:58:18Z"
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "id is required"
}
```

или

```json
{
  "error": "invalid id"
}
```

**Товар не найден (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/items - Получение списка товаров

**URL:** `http://localhost:8080/api/items`

**Authorization:** `Bearer {token}`

**Ожидаемый ответ (200 OK):**

```json
{
  "items": [
    {
      "id": "7a987f51-4841-4d88-ab56-03f3c2fd7d1e",
      "name": "Видеокарта",
      "description": "Palit GeForce RTX 5090 GameRock OC",
      "quantity": 100,
      "price": "313999.00",
      "created_at": "2025-12-24T18:58:18Z",
      "updated_at": "2025-12-24T18:58:18Z"
    }
  ],
  "total": 1
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## PUT /api/items/{id} - Обновление товара

**URL:** `http://localhost:8080/api/items/{id}`

**Content-Type:** `application/json`

**Authorization:** `Bearer {token}` (требует роль admin или manager)

**Параметры:**

- `{id}` (обязательно) - UUID товара
- `name` (опционально) - название товара (минимум 1 символ)
- `description` (опционально) - описание товара
- `quantity` (опционально) - количество товара (минимум 0)
- `price` (опционально) - цена в копейках (минимум 0, максимум 2147483647)

**Body:**

```json
{
  "name": "Видеокарта",
  "description": "Palit GeForce RTX 5090 GameRock OC",
  "quantity": 99,
  "price": 3600900
}
```

**Ожидаемый ответ (200 OK):**

```json
{
  "id": "8204a1b2-3739-49c9-9b59-4d04a1cdb378",
  "name": "Видеокарта",
  "description": "Palit GeForce RTX 5090 GameRock OC",
  "quantity": 99,
  "price": "36009.00",
  "created_at": "2025-12-24T19:05:49Z",
  "updated_at": "2025-12-24T19:05:49Z",
  "message": "item updated successfully"
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "id is required"
}
```

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Quantity' failed on the 'min' tag"
}
```

**Товар не найден (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Доступ запрещен (403 Forbidden):**

```json
{
  "error": "forbidden"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## DELETE /api/items/{id} - Удаление товара

**URL:** `http://localhost:8080/api/items/{id}`

**Authorization:** `Bearer {token}` (требует роль admin)

**Параметры:**

- `{id}` (обязательно) - UUID товара

**Ожидаемый ответ (200 OK):**

```json
{
  "message": "item deleted successfully"
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "id is required"
}
```

**Товар не найден (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Доступ запрещен (403 Forbidden):**

```json
{
  "error": "forbidden"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/items/{id}/history - Получение истории изменений товара

**URL:** `http://localhost:8080/api/items/{id}/history`

**Authorization:** `Bearer {token}`

**Параметры:**

- `{id}` (обязательно) - UUID товара

**Ожидаемый ответ (200 OK):**

```json
{
  "history": [
    {
      "id": "b81beb77-1f53-48eb-82ca-f614ee2d3fd7",
      "item_id": "3f60f58c-de48-4990-9ba0-17a3f9684b4c",
      "action": "create",
      "user_id": "314f393a-6914-457a-98c2-65e89948b198",
      "changed_at": "2025-12-24T18:51:02Z",
      "new_data": {
        "created_at": "2025-12-24T18:51:02.517646+00:00",
        "description": "Palit GeForce RTX 5090 GameRock OC",
        "id": "3f60f58c-de48-4990-9ba0-17a3f9684b4c",
        "name": "Видеокарта",
        "price": 31399900,
        "quantity": 100,
        "updated_at": "2025-12-24T18:51:02.517646+00:00"
      },
      "diff": [
        {
          "field": "created_at",
          "old_value": null,
          "new_value": "2025-12-24T18:51:02.517646+00:00"
        },
        {
          "field": "description",
          "old_value": null,
          "new_value": "Palit GeForce RTX 5090 GameRock OC"
        },
        {
          "field": "id",
          "old_value": null,
          "new_value": "3f60f58c-de48-4990-9ba0-17a3f9684b4c"
        },
        {
          "field": "name",
          "old_value": null,
          "new_value": "Видеокарта"
        },
        {
          "field": "price",
          "old_value": null,
          "new_value": 31399900
        },
        {
          "field": "quantity",
          "old_value": null,
          "new_value": 100
        },
        {
          "field": "updated_at",
          "old_value": null,
          "new_value": "2025-12-24T18:51:02.517646+00:00"
        }
      ]
    }
  ],
  "total": 1
}
```

### Ошибки:

**Некорректный ID (400 Bad Request):**

```json
{
  "error": "id is required"
}
```

**Товар не найден (404 Not Found):**

```json
{
  "error": "item not found"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/history - Получение истории с фильтрами

**URL:** `http://localhost:8080/api/history`

**Authorization:** `Bearer {token}`

**Параметры:**

- `item_id` (опционально) - фильтр по ID товара (UUID)
- `user_id` (опционально) - фильтр по ID пользователя (UUID)
- `action` (опционально) - фильтр по действию: "create", "update", "delete"
- `from` (опционально) - фильтр по дате начала (RFC3339)
- `to` (опционально) - фильтр по дате окончания (RFC3339)
- `sort_by` (опционально) - сортировка: "changed_at", "action", "user_id"
- `sort_order` (опционально) - порядок сортировки: "asc" или "desc"

**Пример запроса:**

```
GET /api/history?item_id=b9ab5b36-444a-47c4-b7b1-7067a4977e67&action=update&from=2025-12-09T00:00:00Z&to=2025-12-09T23:59:59Z&sort_by=changed_at&sort_order=desc
```

**Ожидаемый ответ (200 OK):**

```json
{
  "history": [
    {
      "id": "677b5057-a235-4c1e-8e28-0dc96ea22179",
      "item_id": "8204a1b2-3739-49c9-9b59-4d04a1cdb378",
      "action": "create",
      "user_id": "314f393a-6914-457a-98c2-65e89948b198",
      "changed_at": "2025-12-24T19:05:49Z",
      "new_data": {
        "created_at": "2025-12-24T19:05:49.27508+00:00",
        "description": "Palit GeForce RTX 5090 GameRock OC",
        "id": "8204a1b2-3739-49c9-9b59-4d04a1cdb378",
        "name": "Видеокарта",
        "price": 3600900,
        "quantity": 99,
        "updated_at": "2025-12-24T19:05:49.27508+00:00"
      }
    }
  ],
  "total": 1
}
```

### Ошибки:

**Некорректный формат даты (400 Bad Request):**

```json
{
  "error": "invalid date format"
}
```

```json
{
  "error": "parameter 'from' cannot be after 'to'"
}
```

**Некорректный UUID (400 Bad Request):**

```json
{
  "error": "invalid item_id: invalid UUID format"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Action' failed on the 'action_type' tag"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/history/export - Экспорт истории в CSV

**URL:** `http://localhost:8080/api/history/export`

**Authorization:** `Bearer {token}`

**Параметры:**

- `item_id` (опционально) - фильтр по ID товара (UUID)
- `user_id` (опционально) - фильтр по ID пользователя (UUID)
- `action` (опционально) - фильтр по действию: "create", "update", "delete"
- `from` (опционально) - фильтр по дате начала (RFC3339)
- `to` (опционально) - фильтр по дате окончания (RFC3339)
- `sort_by` (опционально) - сортировка: "changed_at", "action", "user_id"
- `sort_order` (опционально) - порядок сортировки: "asc" или "desc"

**Пример запроса:**

```
GET /api/history/export?from=2025-12-09T00:00:00Z&to=2025-12-09T23:59:59Z&action=update&sort_by=changed_at&sort_order=desc
```

**Ожидаемый ответ (200 OK):**

Файл CSV с заголовками и данными:

```csv
id,item_id,action,user_id,changed_at,old_data,new_data
b2c3d4e5-f6a7-8901-bcde-f12345678901,b9ab5b36-444a-47c4-b7b1-7067a4977e67,update,550e8400-e29b-41d4-a716-446655440000,2025-12-09T20:15:30Z,"{""quantity"":10,""price"":15000000}","{""quantity"":15,""price"":16000000}"
```

**Content-Type:** `text/csv`

**Content-Disposition:** `attachment; filename=history.csv`

### Ошибки:

**Некорректный формат даты (400 Bad Request):**

```json
{
  "error": "invalid date format"
}
```

**Некорректный UUID (400 Bad Request):**

```json
{
  "error": "invalid item_id: invalid UUID format"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "validation for 'Action' failed on the 'action_type' tag"
}
```

**Неавторизован (401 Unauthorized):**

```json
{
  "error": "unauthorized"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```
