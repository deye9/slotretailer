USE retail;

CREATE TABLE IF NOT EXISTS products (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    itemcode     VARCHAR(255) NOT NULL UNIQUE,
	itemname     VARCHAR(255) NOT NULL UNIQUE,
	codebars     VARCHAR(255) NOT NULL,
	onhand       INT,
	minlevel     INT,
	warehouse    VARCHAR(255) NOT NULL,
	serialnumber VARCHAR(255) NOT NULL,
	manufacturer VARCHAR(255) NOT NULL,
	price        REAL,
	vat          REAL
);

CREATE TABLE IF NOT EXISTS customers (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    cardcode    VARCHAR(255) NOT NULL,
	cardname    VARCHAR(255) NOT NULL,
	address     VARCHAR(255) NOT NULL,
	phone       VARCHAR(255) NOT NULL UNIQUE,
    phone1      VARCHAR(255) NOT NULL UNIQUE,
	city        VARCHAR(255) NOT NULL,
	email       VARCHAR(255) NOT NULL,
	synced      INT,
    created_by  INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id        INT AUTO_INCREMENT PRIMARY KEY,
    docentry  INT,
	docnum    INT,
	canceled  INT,
	cardcode  VARCHAR(255) NOT NULL,
	cardname  VARCHAR(255) NOT NULL,
	vatsum    REAL,
	doctotal  REAL,
	synced    INT,
    created_by   INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS ordereditems (
   id       INT AUTO_INCREMENT PRIMARY KEY,
   orderid  INT,
   itemcode VARCHAR(255) NOT NULL,
   itemname VARCHAR(255) NOT NULL,
   price    REAL,
   quantity INT NOT NULL,
   discount REAL NOT NULL DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS users (
   id           INT AUTO_INCREMENT PRIMARY KEY,
   firstname    VARCHAR(255) NOT NULL,
   lastname     VARCHAR(255) NOT NULL,
   email        VARCHAR(255) NOT NULL,
   password     VARCHAR(255) NOT NULL,
   isadmin      INT NOT NULL DEFAULT 0,
   created_by   INT NOT NULL,
   created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at   TIMESTAMP NULL,
   deleted_at   TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS store (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
	address     VARCHAR(255) NOT NULL,
	phone       VARCHAR(255) NOT NULL,
	city        VARCHAR(255) NOT NULL,
	email       VARCHAR(255) NOT NULL,
    orders      VARCHAR(255) NOT NULL,
    products    VARCHAR(255) NOT NULL,
    customers   VARCHAR(255) NOT NULL,
    sync_interval INT NOT NULL DEFAULT 5,
    created_by   INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);

INSERT INTO `users` (firstname, lastname, email, password, created_by, isadmin) 
SELECT 'super', 'admin', 'superadmin@slot.com', 'superadmin', 1, 1
WHERE NOT EXISTS(SELECT * FROM `users` WHERE email = 'superadmin@slot.com' AND password = 'superadmin');
