-- Создаем базу данных, если она не существует
CREATE DATABASE IF NOT EXISTS sorting_app;

-- Используем созданную базу данных
USE sorting_app;

-- Создаем таблицу для хранения массивов
CREATE TABLE IF NOT EXISTS arrays (
    id INT AUTO_INCREMENT PRIMARY KEY,
    array_data TEXT NOT NULL,
    is_sorted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- Создаем пользователя для приложения (замените 'password' на реальный пароль)
CREATE USER IF NOT EXISTS 'sorting_user'@'localhost' IDENTIFIED BY '123';

-- Даем права пользователю на базу данных
GRANT ALL PRIVILEGES ON sorting_app.* TO 'sorting_user'@'localhost';

-- Применяем изменения прав
FLUSH PRIVILEGES;


-- mysql -u sorting_user -p
-- USE sorting_app;
-- SHOW TABLES;
-- DESCRIBE arrays; - структура таблицы
-- SELECT * FROM arrays; - данные в таблице
