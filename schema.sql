-- Schema used by the API.
-- This file is mounted into the MySQL container on first start, so the table
-- is created automatically the first time the stack boots.

CREATE TABLE IF NOT EXISTS products (
    id    INT AUTO_INCREMENT PRIMARY KEY,
    name  VARCHAR(255)   NOT NULL,
    stock INT            NOT NULL DEFAULT 0,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0
);
