-- Создание таблицы users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL
);

-- Создание таблицы delivery_items
CREATE TABLE delivery_items (
    id SERIAL PRIMARY KEY,
    image VARCHAR(255),
    title VARCHAR(255),
    price BIGINT NOT NULL,
    description TEXT,
    is_delete BOOLEAN DEFAULT false
);


-- Создание таблицы delivery_requests
CREATE TABLE delivery_requests (
    id SERIAL PRIMARY KEY,
    status VARCHAR(255),
    address VARCHAR(255),
    delivery_date TIMESTAMP,
    delivery_type VARCHAR(255),
    user_id BIGINT,
    moderator_id BIGINT
);

-- Создание таблицы item_requests
CREATE TABLE item_requests (
    id PRIMARY KEY,
    item_id BIGINT NOT NULL,
    request_id BIGINT NOT NULL,
    count BIGINT NOT NULL DEFAULT 1,
    PRIMARY KEY (item_id, request_id)
);
