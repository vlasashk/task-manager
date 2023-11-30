# Task-manager

- POST /api/task - Создание задачи
    ```
    body
    {
        "title": "Название задачи",
        "description": "Описание задачи",
        "due_date": "Дата завершения задачи",
        "status": "Выполнено/Не выполнено"
    }
    ```
- GET /api/tasks - Получение списка задач
    > Параметры запроса:
  > - status (bool, optional): Фильтр по статусу задачи (true - выполнено, false - не выполнено). По дефолту выводит оба типа.
  > - date (string, optional): Фильтр по дате задачи. (Формат YYYY-MM-DD).
  > - page (uint, optional): Номер страницы. 

  >GET /tasks?status=true&date=2023-11-29&page=0
- GET /api/task/{id} - Получение задачи по ID
- PUT /api/task/{id} - Обновление задачи
    ```
    body
    {
        "title": "Название задачи",
        "description": "Описание задачи",
        "due_date": "Дата завершения задачи",
        "status": "Выполнено/Не выполнено"
    }
    ```
- DELETE /api/task/{id} - Удаление задачи