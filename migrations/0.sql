CREATE DATABASE IF NOT EXISTS shopdb;

USE shopdb;

CREATE TABLE requests (
    id INT(6) AUTO_INCREMENT PRIMARY KEY,
    payment_id VARCHAR(100),
    item_name VARCHAR(100),
    price FLOAT,
    tax FLOAT,
    track VARCHAR(100),
    status VARCHAR(100),
    created_at datetime,
    store_id INT,
    product_id INT
)