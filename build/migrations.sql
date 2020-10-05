USE retail;

CREATE TABLE IF NOT EXISTS products (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    itemcode     VARCHAR(255) NOT NULL UNIQUE,
	itemname     VARCHAR(255) NOT NULL,
	codebars     VARCHAR(255) NOT NULL,
	onhand       INT,
	minlevel     INT,
	warehouse    VARCHAR(255) NOT NULL,
	serialnumber VARCHAR(255) NOT NULL,
	manufacturer VARCHAR(255) NOT NULL,
	price        REAL,
	vat          REAL
);

CREATE TABLE IF NOT EXISTS banks (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(255) NOT NULL UNIQUE,
    code         VARCHAR(255) NOT NULL
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
	synced      boolean,
    created_by  INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id        INT AUTO_INCREMENT PRIMARY KEY,
    docentry  INT,
	docnum    INT,
	canceled  boolean,
	cardcode  VARCHAR(255) NOT NULL,
	cardname  VARCHAR(255) NOT NULL,
	vatsum    REAL,
	doctotal  REAL,
	synced    boolean,
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

CREATE TABLE IF NOT EXISTS payments (
   id       INT AUTO_INCREMENT PRIMARY KEY,
   orderid  INT,
   docentry INT,
   docnum   INT,
   canceled  boolean,
   paymenttype VARCHAR(255) NOT NULL,
   paymentdetails VARCHAR(255) NOT NULL DEFAULT "",
   amount   REAL  NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
   id           INT AUTO_INCREMENT PRIMARY KEY,
   firstname    VARCHAR(255) NOT NULL,
   lastname     VARCHAR(255) NOT NULL,
   email        VARCHAR(255) NOT NULL,
   password     VARCHAR(255) NOT NULL,
   isadmin      boolean NOT NULL DEFAULT false,
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
    banks       VARCHAR(255) NOT NULL,
    sync_interval INT NOT NULL DEFAULT 5,
    sapkey      VARCHAR(255) NOT NULL,
    created_by   INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);

IF NOT EXISTS( SELECT NULL
            FROM INFORMATION_SCHEMA.COLUMNS
           WHERE table_name = 'orders'
             AND table_schema = 'retail'
             AND column_name = 'comment')  THEN

ALTER table orders add column comment text;
ALTER table orders add column returned boolean DEFAULT false;

END IF;

REPLACE INTO `users` (firstname, lastname, email, password, created_by, isadmin) 
SELECT 'super', 'admin', 'superadmin@slot.com', 'superadmin', 1, true
WHERE NOT EXISTS(SELECT * FROM `users` WHERE email = 'superadmin@slot.com' AND password = 'superadmin');

CREATE PROCEDURE `searchDB`(tableSchema nvarchar(50), searchTerm nvarchar(50))
BEGIN
SET @condition = concat("LIKE '%", searchTerm, "%'");
SET @column_types_regexp = '^((var)?char|(var)?binary|blob|text|enum|set)\\(';

-- Reset @sql_query in case it was used previously
SET @sql_query = '';

-- Build query for each table and merge with previous queries with UNION
SELECT
    -- Use 'DISTINCT IF(QUERYBUILDING, NULL, NULL)'
    -- to only select a single null value
    -- instead of selecting the query over and over again as it's built
    DISTINCT IF(@sql_query := CONCAT(
        IF(LENGTH(@sql_query), CONCAT(@sql_query, " UNION "), ""),
        'SELECT ',
            QUOTE(CONCAT('`', `table_name`, '`.`', `column_name`, '`')), ' AS `column`, ',
            'COUNT(*) AS occurrences ',
        'FROM `', `table_schema`, '`.`', `table_name`, '` ',
        'WHERE `', `column_name`, '` ', @condition
    ), NULL, NULL) query
FROM (
    SELECT table_schema, table_name, column_name FROM information_schema.columns
    WHERE table_schema = tableSchema AND column_type REGEXP @column_types_regexp
) results;

-- Only return results with at least one occurrence
SET @sql_query = CONCAT("SELECT * FROM (", @sql_query, ") results WHERE occurrences > 0");

-- Run built query
PREPARE statement FROM @sql_query;
EXECUTE statement;
DEALLOCATE PREPARE statement;

END

-- ALTER table orders add column comment text;