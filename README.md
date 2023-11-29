# TODO-manager

- POST /api/tasks - Создание задачи
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
  > - status (boolean, optional): Фильтр по статусу задачи (true - выполнено, false - не выполнено). По дефолту выводит оба типа. 
  > - page (integer, optional): Номер страницы. 
  > - limit (integer, optional): Количество задач на странице. 

  >GET /tasks?status=true&page=1&limit=10
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