CREATE TABLE IF NOT EXISTS tasks (
     id VARCHAR(255) PRIMARY KEY,
     title VARCHAR(255) NOT NUll,
     description TEXT NOT NUll,
     due_date DATE NOT NULL CHECK (due_date >= CURRENT_DATE),
     status BOOLEAN DEFAULT FALSE,
     deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tasks_not_deleted ON tasks (id) WHERE deleted_at IS NULL;
