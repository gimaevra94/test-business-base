# Система управления сервисными заявками

Веб-приложение для управления заявками на ремонт с разграничением прав доступа по ролям.

## 📋 Возможности

- **Роли пользователей**: Диспетчер и Мастер
- **Управление заявками**: Создание, назначение, отслеживание статуса
- **Жизненный цикл заявки**: new → assigned → in_progress → done/canceled
- **Панель управления**: Фильтрация заявок по статусу
- **Оптимистичная блокировка**: Контроль версий для предотвращения гонок данных

## 🏗️ Архитектура проекта

```
├── consts/          # Константы (SQL-запросы, пути, сообщения)
├── database/        # Подключение к БД и запросы
├── handlers/        # Обработчики HTTP-запросов
├── structs/         # Структуры данных
├── templates/       # HTML-шаблоны
├── migrations/      # Миграции базы данных
├── seeds/           # Начальные данные
└── docker/          # Docker-конфигурация
```

## 🚀 Быстрый старт

### Требования
- Docker и Docker Compose
- Go 1.25+ (для локальной разработки)

### Запуск через Docker Compose

```bash
cd docker
docker-compose up --build
```

Приложение будет доступно по адресам:
- **Веб-интерфейс**: http://localhost:8080
- **MySQL**: localhost:3306

### Локальная разработка

1. **Запуск базы данных**:
```bash
cd docker
docker-compose up db
```

2. **Применение миграций**:
```bash
docker-compose run --rm migrate -path ../migrations -database "mysql://root:root@tcp(localhost:3306)/repair_service" up
```

3. **Заполнение тестовыми данными**:
```bash
docker-compose run --rm seed
```

4. **Запуск приложения**:
```bash
cd app
go run main.go
```

## 👥 Роли пользователей

### Диспетчер
- Создание новых заявок
- Назначение мастеров на заявки (статус: new)
- Отмена заявок (статус: new, assigned)
- Просмотр всех заявок

### Мастер
- Просмотр назначенных заявок
- Начало работы (статус: assigned → in_progress)
- Завершение работы (статус: in_progress → done)

## 🗄️ Схема базы данных

**Таблица users**
- id, name, role (dispatcher/master)

**Таблица requests**
- id, client_name, phone, address, problem_text
- status (new/assigned/in_progress/done/canceled)
- assigned_to (внешний ключ → users.id)
- version (оптимистичная блокировка)
- created_at, updated_at

## 🔧 Конфигурация

Подключение к базе данных:
```
Host: localhost:3306
User: root
Password: root
Database: repair_service
```

## 📦 Технологический стек

- **Backend**: Go 1.25, маршрутизатор chi
- **База данных**: MySQL 8.0
- **Работа с БД**: GORM + raw SQL
- **Шаблоны**: HTML templates
- **Контейнеризация**: Docker, Docker Compose

## 🔐 Аутентификация

Простая аутентификация на основе cookie с выбором пользователя при входе.

## 📝 Эндпоинты API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/` | Страница входа |
| POST | `/login` | Аутентификация |
| GET | `/logout` | Выход из системы |
| GET/POST | `/create` | Создание заявки |
| GET | `/dashboard` | Просмотр заявок |
| POST | `/action` | Выполнение действий (assign/cancel/start/finish) |

## 🧪 Тестовые пользователи

После заполнения seed-данными доступны:
- **Диспетчеры**: Dispatcher 2, Dispatcher 3
- **Мастера**: Master 3, Master 4, Master 5, Master 6

## 📄 Лицензия

MIT