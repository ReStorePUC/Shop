USE shopdb;

CREATE TABLE payments (
      id INT(6) AUTO_INCREMENT PRIMARY KEY,
      total FLOAT,
      pix VARCHAR(100),
      status VARCHAR(100),
      created_at datetime,
      store_id INT,
      product_id INT
)
