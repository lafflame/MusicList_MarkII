# MusicList MarkII

## Описание
REST API для управления личной музыкальной коллекцией. Поддерживает треки, плейлисты, поиск и статистику.

## Стек
- **Go** + Gin
- **PostgreSQL** + GORM
- **Docker** + Docker Compose
- **React** (фронтенд)
- **Swagger** (документация API)

## Структура проекта
```
MusicList_MarkII/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/        # конфигурация БД
│   ├── domain/        # модели
│   ├── repository/    # работа с БД
│   ├── service/       # бизнес-логика
│   └── handler/       # HTTP хендлеры
├── frontend/          # React приложение
├── docs/              # Swagger
├── Dockerfile
├── docker-compose.yml
└── .env
```

## Запуск через Docker

1. Склонируй репозиторий:
```bash
git clone <url>
cd MusicList_MarkII
```

2. Создай `.env` в корне:
```
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=mydb
DB_PORT=5432
```

3. Запусти:
```bash
docker-compose up --build
```

После запуска:
- Фронтенд: `http://localhost:3000`
- Бэкенд: `http://localhost:8081`
- Swagger: `http://localhost:8081/swagger/index.html`

## Локальный запуск (без Docker)

### Зависимости
- Go 1.25+
- PostgreSQL

### Установка
```bash
go mod download
```

### Настройка БД
Создай базу данных `mydb` в PostgreSQL и заполни `.env`:
```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=mydb
DB_PORT=5432
```

### Запуск
```bash
go run ./cmd/
```

## API

### Треки
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/tracks` | Все треки |
| POST | `/api/tracks` | Добавить трек |
| PUT | `/api/tracks/:id` | Обновить трек |
| DELETE | `/api/tracks/:id` | Удалить трек |
| GET | `/api/tracks/search?query=` | Поиск |
| GET | `/api/tracks/filter?from=&to=` | Фильтр по дате |
| GET | `/api/tracks/shuffle` | Перемешать |

### Плейлисты
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/playlists` | Все плейлисты |
| POST | `/api/playlists` | Создать плейлист |
| PUT | `/api/playlists/:id` | Переименовать |
| DELETE | `/api/playlists/:id` | Удалить |
| GET | `/api/playlists/:id/tracks` | Треки плейлиста |
| POST | `/api/playlists/:id/tracks/:track_id` | Добавить трек в плейлист |
| DELETE | `/api/playlists/:id/tracks/:track_id` | Убрать трек из плейлиста |

### Статистика
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/statistics` | Статистика медиатеки |