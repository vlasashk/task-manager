# Task-manager
## Build
#### Prerequisites
- docker

1. Clone project:
```
git clone https://github.com/vlasashk/task-manager.git
cd task-manager
```
2. Run:
```
docker compose up --build
```
3. Test:
```
go test -v ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html
```
## Project information

### Restrictions
- Дата должа передаваться в валидном формате (YYYY-MM-DD)
- Дата на которую заводится задача не может быть ранее текущей даты
- Лимит на вывод списка задач на 1 страницу захардкожен на значении = 10 записей
- Пагинация - единственный режим взаимодействия с API (нельзя получить более 10 записей за 1 запрос)
- Пагинация начинается со странцы = 0
- Что бы получить список задач - не обязательно передавать параметры (page, date, status), тогда будут выведены первые 10 задач(не зависимо от статуса), отсортированные по дате
- Вывод списка задач с фильтрацией по дате(без фильтраци по статусу) - выведет все задачи аткуальные на конкретную дату
- Возможна одовременная фильтрация и по дате и по статусу
- Удаление задачи, не удаляет запись из БД, а помечает как удаленную

### Tools used
- PostgreSQL as database
- [jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx) package as toolkit for PostgreSQL
- [go-chi/chi](https://pkg.go.dev/github.com/go-chi/chi) package as router for building HTTP service
- [swaggo/swag](https://github.com/swaggo/swag) package as swagger doc generator
- [rs/zerolog](https://github.com/rs/zerolog) package for logging
- [stretchr/testify](https://github.com/stretchr/testify) package for testing
- [vektra/mockery](https://github.com/vektra/mockery) package for mock generation
- Docker for deployment

### Functionality
#### Swagger
Swagger generated documentation will be available after run at `http://localhost:9090/api/swagger/index.html` (or different port if .env file was edited)

#### Tasks manipulation
- {POST} /api/task - Создание задачи
    ```
    body
    {
        "title": "Название задачи",
        "description": "Описание задачи",
        "due_date": "Дата завершения задачи",
        "status": "Выполнено/Не выполнено"
    }
    ```
- {GET} /api/tasks - Получение списка задач
    > Параметры запроса:
  > - status (bool, optional): Фильтр по статусу задачи (true - выполнено, false - не выполнено). По дефолту выводит оба типа.
  > - date (string, optional): Фильтр по дате задачи. (Формат YYYY-MM-DD).
  > - page (uint, optional): Номер страницы. 

  > {GET} /api/tasks?status=false&date=2024-12-29&page=0
- {GET} /api/task/{id} - Получение задачи по ID
- {PUT} /api/task/{id} - Обновление задачи
    ```
    body
    {
        "title": "Название задачи",
        "description": "Описание задачи",
        "due_date": "Дата завершения задачи",
        "status": "Выполнено/Не выполнено"
    }
    ```
- {DELETE} /api/task/{id} - Удаление задачи