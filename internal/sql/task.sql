CREATE TABLE IF NOT EXISTS tasks (
     id VARCHAR(255) PRIMARY KEY,
     title VARCHAR(255) PRIMARY KEY,
     description TEXT NOT NUll,
     due_date DATE NOT NULL CHECK (due_date >= CURRENT_DATE),
     completed BOOLEAN DEFAULT FALSE
);