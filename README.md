# 📦 Project Manager Dashboard (Go)

Backend-сервис для управления проектами и задачами.  
Реализован на **Go**, **Ent ORM**, **PostgreSQL**, **Chi Router**.

**Frontend:** [project-manager-dashboard-angular](https://github.com/donuwave/project-manager-dashboard-angular)

![Go](https://img.shields.io/badge/Go_1.24-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white)
![Ent](https://img.shields.io/badge/Ent_ORM-4B32C3?style=flat)
![Chi](https://img.shields.io/badge/Chi_Router-000000?style=flat)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)

---

## 🚀 Возможности

### Пользователи
- Создание пользователя
- Получение пользователя по ID
- Список пользователей

### Проекты
- Создание проекта
- Получение проекта с участниками
- Список проектов
- Приглашение пользователя в проект
- Удаление проекта (**только owner**)

### Задачи
- Создание задачи в проекте
- Список задач проекта
- Обновление задачи (PATCH)
- Назначение задачи пользователю / на себя
- Получение assignee в списке задач
- Удаление задачи (**только owner проекта**)
- Изменение порядка задач (`position`): сдвиг соседних задач выполняется в одной транзакции

---

## 🧱 Архитектура

Проект построен по принципам **Clean Architecture**.

```
cmd/
 ├── api/                 # точка входа HTTP-сервера
 └── migrate/             # создание схемы БД через ent
internal/
 ├── app/
 │   ├── usecase/         # бизнес-логика
 │   │   ├── user/
 │   │   ├── project/
 │   │   └── task/        # service + repository + типы и ошибки
 │   └── app.go           # инициализация ent
 └── transport/
     └── http/            # handlers, router, DTO, мапперы
ent/
 ├── schema/              # схема: user, project, task, project_user, project_task
 └── ...                  # сгенерированный код ent
```

Связи «многие-ко-многим» оформлены отдельными сущностями с собственными полями: `project_user` (участник проекта) и `project_task` (задача в проекте с позицией на доске).

## 📡 API

| Метод | Путь | Описание |
| --- | --- | --- |
| `GET` | `/users` | Список пользователей |
| `POST` | `/users` | Создание пользователя |
| `GET` | `/users/{id}` | Пользователь по ID |
| `GET` | `/projects` | Список проектов |
| `POST` | `/projects` | Создание проекта |
| `GET` | `/projects/{id}` | Проект с участниками |
| `PATCH` | `/projects/{id}` | Изменение названия и описания |
| `POST` | `/projects/{id}/invite` | Приглашение пользователя |
| `DELETE` | `/projects/{id}` | Удаление проекта (только owner) |
| `GET` | `/projects/{id}/tasks` | Задачи проекта |
| `POST` | `/projects/{id}/tasks` | Создание задачи |
| `PATCH` | `/tasks/{id}` | Изменение задачи: поля, статус, позиция |
| `POST` | `/tasks/{id}/assign` | Назначение исполнителя |
| `DELETE` | `/tasks/{id}` | Удаление задачи (только owner проекта) |

## 🛠 Технологии

- Go 1.22+
- PostgreSQL
- Ent ORM
- Chi Router
- UUID
- pgx
- dotenv

---

## Установка зависимостей

```sh
    go mod tidy
    ent generate ./ent/schema
    go run ./cmd/api
```

## Запуск сервиса

```sh
  docker compose up --build
```

## Документация

Коллекция запросов для Insomnia — [`doc.yaml`](doc.yaml) в корне проекта: Insomnia → Import → выбрать файл.
