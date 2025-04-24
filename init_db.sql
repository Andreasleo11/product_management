CREATE DATABASE product_management;

\ c product_management CREATE TABLE users (
  id SERIAL PRIMARY KEY
  , name VARCHAR(100) NOT NULL
  , email VARCHAR(150) NOT NULL
  , password VARCHAR(255) NOT NULL
  , role VARCHAR(50)
  , created_at TIMESTAMP NOT NULL DEFAULT NOW()
  , updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL);

CREATE TABLE products (
  id SERIAL PRIMARY KEY
  , name VARCHAR(100) NOT NULL
  , description TEXT
  , purchase_price NUMERIC
  , sell_price NUMERIC
  , category_id INTEGER REFERENCES categories(id)
  , stock INTEGER
  , min_stock_alert INTEGER
  , image_url TEXT
  , status VARCHAR(50)
  , created_at TIMESTAMP NOT NULL DEFAULT NOW()
  , updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE stock_logs (
  id SERIAL PRIMARY KEY
  , product_id INT REFERENCES products(id)
  , change_type VARCHAR(10)
  , -- 'in' atau 'out'
    amount INT
  , note TEXT
  , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);