drop database retail;

create database retail;
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

CREATE TABLE IF NOT EXISTS pricelist (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(255) NOT NULL UNIQUE,
    code         VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS cashaccounts (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(255) NOT NULL UNIQUE,
    code         VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS stores (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL UNIQUE,
  code varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
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
    -- gender      VARCHAR(255) NOT NULL,
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

CREATE TABLE IF NOT EXISTS reports (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL UNIQUE,
    query       TEXT NOT NULL,
    created_by  INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS audits (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    old_row_data JSON NULL,
    new_row_data JSON NULL,
    dml_type    ENUM('INSERT', 'UPDATE', 'DELETE') NOT NULL,
    dml_created_by  INT NULL DEFAULT 0,
    dml_timestamp  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transfers(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    fromwhs     int,
    towhs       int,
    comment     text,
    canceled    boolean,
    synced      boolean,
    created_by  INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS transfereditems (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    transferid  INT,
    itemcode    VARCHAR(255) NOT NULL,
    itemname    VARCHAR(255) NOT NULL,
    onHand      INT,
    quantity    INT
);

ALTER table orders add column comment text;
ALTER table orders add column returned boolean DEFAULT false;
ALTER table store add column logrotation text NOT NULL;
ALTER table store add column transfers text NULL;
ALTER table store add column vat boolean DEFAULT false;
ALTER table orders add column discountapprovedby int DEFAULT 0;
ALTER table store add column warehouses VARCHAR(255) NOT NULL;
ALTER TABLE products MODIFY serialnumber TEXT NOT NULL;
ALTER TABLE products RENAME column serialnumber to serialnumbers; 
ALTER table ordereditems add column serialnumber VARCHAR(255) NOT NULL;
ALTER table store rename column banks to creditcard;
ALTER table store add column pricelist VARCHAR(255) NOT NULL;
ALTER table store add column productpricelist VARCHAR(255) NOT NULL;
ALTER TABLE banks RENAME TO creditcards;
ALTER table store add column cashaccount VARCHAR(255) NOT NULL;
ALTER table store add column storecashaccount VARCHAR(255) NOT NULL;

-- Default Reports
REPLACE INTO reports (id, title, query, created_by) VALUES (1, "Todays Orders", "select id as order_id, docnum `Document Number`, canceled `Is Cancelled`, CardCode, CardName,  vatsum `VAT %`, concat('₦', format(doctotal, 2)) `Document Total`, case when Synced <> 0 then \"Yes\" else \"No\" END `Synced`, case when returned <> 0 then \"Yes\" else \"No\" END `Returned`, ifnull( (select concat(firstname, '  ', lastname) from users where users.id = o.discountapprovedby), 'Super Admin') `approved_by` from orders o where deleted_at is null and cast(created_at as date) = CURDATE() order by created_at desc;", 1);
REPLACE INTO reports (id, title, query, created_by) VALUES (2, "Week's Top Seller", "select i.id, i.orderid, i.itemcode, i.itemname, sum(i.price) price, sum(i.quantity) quantity, sum(i.discount) discount from ordereditems i inner join orders o on o.id = i.orderid WHERE cast(created_at as date) BETWEEN cast(DATE_ADD(CURDATE(), INTERVAL(1 - DAYOFWEEK(CURDATE())) DAY) as date) AND cast(DATE_ADD(CURDATE(), INTERVAL(7 - DAYOFWEEK(CURDATE())) DAY) as date) AND deleted_At is null GROUP BY i.id, i.orderid, i.itemcode, i.itemname;", 1);
REPLACE INTO reports (id, title, query, created_by) VALUES (3, "Todays Top Sellers", "select i.* from orders o inner join ordereditems i on o.id = i.orderid where deleted_at is null and cast(created_at as date) = CURDATE() order by created_at desc;", 1);
REPLACE INTO reports (id, title, query, created_by) VALUES (4, "Store Inventory Level", "select p.* from store s inner join products p on s.sapkey = p.warehouse;", 1);
REPLACE INTO reports (id, title, query, created_by) VALUES (5, "Global Inventory Level", "select * from products p;", 1);
INSERT INTO store (`id`,`name`,`address`,`phone`,`city`,`email`,`orders`,`products`,`customers`,`creditcard`,`sync_interval`,`sapkey`,`created_by`,`created_at`,`updated_at`,`deleted_at`,`logrotation`,`transfers`,`vat`,`warehouses`, `pricelist`, `productpricelist`, `cashaccount`, `storecashaccount`) VALUES (1,'New Store','Enter Store Address','080','Lagos','storename@slot.com','http://197.255.32.34:5000/Orders','http://197.255.32.34:5000/Products','http://197.255.32.34:5000/Customers','http://197.255.32.34:5000/CreditCards',30,'2BMEDICA',1,'2020-12-08 21:46:36',NULL,NULL,'1','http://197.255.32.34:5000/Transfers',0,'http://197.255.32.34:5000/Warehouses', 'http://197.255.32.34:5000/pricelist', 'code 1', 'http://197.255.32.34:5000/CashAccounts', '12330001');
REPLACE INTO `users` (firstname, lastname, email, password, created_by, isadmin) VALUES ('super', 'admin', 'superadmin@slot.com', 'superadmin', 1, true);

-- REPLACE INTO `users` (firstname, lastname, email, password, created_by, isadmin) SELECT 'super', 'admin', 'superadmin@slot.com', 'superadmin', 1, true WHERE NOT EXISTS(SELECT * FROM `users` WHERE email = 'superadmin@slot.com' AND password = 'superadmin');

-- Get Orders Stored Procedure
USE `retail`;
DROP procedure IF EXISTS `getOrders`;

USE `retail`;
DROP procedure IF EXISTS `retail`.`getOrders`;

DELIMITER $$
USE `retail`$$
CREATE DEFINER=`root`@`localhost` PROCEDURE `getOrders`()
BEGIN
	declare orderID integer;
	declare finished integer default 0;
    
    declare order_cursor cursor for
        select id from orders where synced = false and canceled = false;

	/* declare not found handler*/
	declare continue handler
	for not found set finished = 1;

	 open order_cursor;

	get_details : LOOP
     fetch order_cursor into orderID;

      if finished = 1 THEN
      leave get_details;
      end if;

     /* Get order alongside ordereditems and payments */
     select id, DATE_FORMAT(created_at, '%Y-%m-%d') docdate, docnum,
			case when cardcode like 'R-%' then replace(cardcode, 'R-', '') else 0 end customer_id, 
            case when cardcode like 'R-%' then 0 else replace(cardcode, 'R-', '') end cardcode,
            cardname, doctotal,  comment, returned, synced,
            (select JSON_ARRAYAGG(JSON_OBJECT(
			 'itemcode', itemcode, 
			 'itemname', itemname, 
			 'quantity', quantity, 
			 'price', price, 
			 'discount', discount, 
			'warehouse', (select sapkey from store) ,
			'serialnumber', serialnumber))  
			from ordereditems where orderid = orderID) items,
            (select JSON_ARRAYAGG(JSON_OBJECT(
			 'amount', amount, 
			 'paymenttype', paymenttype,
			 'paymentdetails', case when lower(paymenttype) = 'cash' then (select storecashaccount from store) else paymentdetails end))
			 from payments where orderid = orderID) payments 
             from orders where id = orderID;

     END LOOP get_details;
close order_cursor;
END$$

DELIMITER ;

-- CUSTOMERS TRIGGER
drop trigger if exists customer_insert_audit;
drop trigger if exists customer_update_audit;
drop trigger if exists customer_delete_audit;

DELIMITER $$

CREATE trigger customer_insert_audit
AFTER INSERT ON customers FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (null, 
    JSON_OBJECT('id', new.id, 'cardcode', new.cardcode, 'cardname', new.cardname, 'address', new.address, 'phone', new.phone, 'phone1', new.phone1, 'city', new.city, 'email', new.email, 'synced', new.synced, 'created_at', new.created_at, 'created_by', new.created_by), 
    'INSERT', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger customer_update_audit
AFTER UPDATE ON customers FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'cardcode', old.cardcode, 'cardname', old.cardname, 'address', old.address, 'phone', old.phone, 'phone1', old.phone1, 'city', old.city, 'email', old.email, 'synced', old.synced, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
    JSON_OBJECT('id', new.id, 'cardcode', new.cardcode, 'cardname', new.cardname, 'address', new.address, 'phone', new.phone, 'phone1', new.phone1, 'city', new.city, 'email', new.email, 'synced', new.synced, 'created_at', new.created_at, 'updated_at', new.updated_at, 'deleted_at', new.deleted_at, 'created_by', new.created_by), 
    'UPDATE', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger customer_delete_audit
AFTER DELETE ON customers FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'cardcode', old.cardcode, 'cardname', old.cardname, 'address', old.address, 'phone', old.phone, 'phone1', old.phone1, 'city', old.city, 'email', old.email, 'synced', old.synced, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
	null, 'DELETE', old.created_by);
END $$    

DELIMITER ;

-- USERS TRIGGER
drop trigger if exists user_insert_audit;
drop trigger if exists user_update_audit;
drop trigger if exists user_delete_audit;

DELIMITER $$

CREATE trigger user_insert_audit
AFTER INSERT ON users FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (null, 
    JSON_OBJECT('id', new.id, 'firstname', new.firstname, 'lastname', new.lastname, 'email', new.email, 'password', new.password, 'isadmin', new.isadmin, 'created_at', new.created_at, 'created_by', new.created_by), 
    'INSERT', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger user_update_audit
AFTER UPDATE ON users FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'firstname', old.firstname, 'lastname', old.lastname, 'email', old.email, 'password', old.password, 'isadmin', old.isadmin, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
    JSON_OBJECT('id', new.id, 'firstname', new.firstname, 'lastname', new.lastname, 'email', new.email, 'password', new.password, 'isadmin', new.isadmin, 'created_at', new.created_at, 'updated_at', new.updated_at, 'deleted_at', new.deleted_at, 'created_by', new.created_by), 
    'UPDATE', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger user_delete_audit
AFTER DELETE ON users FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'firstname', old.firstname, 'lastname', old.lastname, 'email', old.email, 'password', old.password, 'isadmin', old.isadmin, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
	null, 'DELETE', old.created_by);
END $$    

DELIMITER ;

-- STORES TRIGGER
drop trigger if exists store_insert_audit;
drop trigger if exists store_update_audit;
drop trigger if exists store_delete_audit;

DELIMITER $$

CREATE trigger store_insert_audit
AFTER INSERT ON store FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (null, 
    JSON_OBJECT('id', new.id, 'name', new.name, 'address', new.address, 'phone', new.phone, 'city', new.city, 'email', new.email, 'orders', new.orders, 'products', new.products, 'customers', new.customers, 'banks', new.banks, 'sync_interval', new.sync_interval, 'sapkey', new.sapkey, 'logrotation', new.logrotation, 'created_at', new.created_at, 'created_by', new.created_by), 
    'INSERT', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger store_update_audit
AFTER UPDATE ON store FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'name', old.name, 'address', old.address, 'phone', old.phone, 'city', old.city, 'email', old.email, 'orders', old.orders, 'products', old.products, 'customers', old.customers, 'banks', old.banks, 'sync_interval', old.sync_interval, 'sapkey', old.sapkey, 'logrotation', old.logrotation, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
    JSON_OBJECT('id', new.id, 'name', new.name, 'address', new.address, 'phone', new.phone, 'city', new.city, 'email', new.email, 'orders', new.orders, 'products', new.products, 'customers', new.customers, 'banks', new.banks, 'sync_interval', new.sync_interval, 'sapkey', new.sapkey, 'logrotation', new.logrotation, 'created_at', new.created_at, 'updated_at', new.updated_at, 'deleted_at', new.deleted_at, 'created_by', new.created_by), 
    'UPDATE', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger store_delete_audit
AFTER DELETE ON store FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'name', old.name, 'address', old.address, 'phone', old.phone, 'city', old.city, 'email', old.email, 'orders', old.orders, 'products', old.products, 'customers', old.customers, 'banks', old.banks, 'sync_interval', old.sync_interval, 'sapkey', old.sapkey, 'logrotation', old.logrotation, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
	null, 'DELETE', old.created_by);
END $$    

DELIMITER ;

-- ORDERS TRIGGER
drop trigger if exists order_insert_audit;
drop trigger if exists order_update_audit;
drop trigger if exists order_delete_audit;

DELIMITER $$

CREATE trigger order_insert_audit
AFTER INSERT ON orders FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (null, 
    JSON_OBJECT('id', new.id, 'docentry', new.docentry, 'docnum', new.docnum, 'canceled', new.canceled, 'cardcode', new.cardcode, 'cardname', new.cardname, 'vatsum', new.vatsum, 'doctotal', new.doctotal, 'synced', new.synced, 'comment', new.comment, 'returned', new.returned, 'created_at', new.created_at, 'created_by', new.created_by), 
    'INSERT', new.created_by);
END $$

DELIMITER ;

DELIMITER $$

CREATE trigger order_update_audit
AFTER UPDATE ON orders FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'docentry', old.docentry, 'docnum', old.docnum, 'canceled', old.canceled, 'cardcode', old.cardcode, 'cardname', old.cardname, 'vatsum', old.vatsum, 'doctotal', old.doctotal, 'synced', old.synced, 'comment', old.comment, 'returned', old.returned, 'created_by', old.created_by, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at),
    JSON_OBJECT('id', new.id, 'docentry', new.docentry, 'docnum', new.docnum, 'canceled', new.canceled, 'cardcode', new.cardcode, 'cardname', new.cardname, 'vatsum', new.vatsum, 'doctotal', new.doctotal, 'synced', new.synced, 'comment', new.comment, 'returned', new.returned, 'created_by', new.created_by, 'created_at', new.created_at, 'updated_at', new.updated_at, 'deleted_at', new.deleted_at), 
    'UPDATE', new.created_by);
END $$    

DELIMITER ;

DELIMITER $$

CREATE trigger order_delete_audit
AFTER DELETE ON orders FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'docentry', old.docentry, 'docnum', old.docnum, 'canceled', old.canceled, 'cardcode', old.cardcode, 'cardname', old.cardname, 'vatsum', old.vatsum, 'doctotal', old.doctotal, 'synced', old.synced, 'comment', old.comment, 'returned', old.returned, 'created_at', old.created_at, 'updated_at', old.updated_at, 'deleted_at', old.deleted_at, 'created_by', old.created_by),
	null, 'DELETE', old.created_by);
END $$    

DELIMITER ;

-- ORDEREDITEMS TRIGGER
drop trigger if exists ordereditems_insert_audit;
drop trigger if exists ordereditems_update_audit;
drop trigger if exists ordereditems_delete_audit;

DELIMITER $$

CREATE trigger ordereditems_insert_audit
AFTER INSERT ON ordereditems FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (null, 
    JSON_OBJECT('id', new.id, 'orderid', new.orderid, 'itemcode', new.itemcode, 'itemname', new.itemname, 'price', new.price, 'quantity', new.quantity, 'discount', new.discount), 
    'INSERT', null);
END $$

DELIMITER ;

DELIMITER $$

CREATE trigger ordereditems_update_audit
AFTER UPDATE ON ordereditems FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'orderid', old.orderid, 'itemcode', old.itemcode, 'itemname', old.itemname, 'price', old.price, 'quantity', old.quantity, 'discount', old.discount),
    JSON_OBJECT('id', new.id, 'orderid', new.orderid, 'itemcode', new.itemcode, 'itemname', new.itemname, 'price', new.price, 'quantity', new.quantity, 'discount', new.discount), 
    'UPDATE', null);
END $$

DELIMITER ;

DELIMITER $$

CREATE trigger ordereditems_delete_audit
AFTER DELETE ON ordereditems FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'orderid', old.orderid, 'itemcode', old.itemcode, 'itemname', old.itemname, 'price', old.price, 'quantity', old.quantity, 'discount', old.discount),
	null, 'DELETE', null);
END $$

DELIMITER ;

-- PAYMENTS TRIGGER
drop trigger if exists payments_insert_audit;
drop trigger if exists payments_update_audit;
drop trigger if exists payments_delete_audit;

DELIMITER $$

CREATE trigger payments_insert_audit
AFTER INSERT ON payments FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (null, 
    JSON_OBJECT('id', new.id, 'orderid', new.orderid, 'docentry', new.docentry, 'docnum', new.docnum, 'canceled', new.canceled, 'paymenttype', new.paymenttype, 'paymentdetails', new.paymentdetails, 'amount', new.amount), 
    'INSERT', null);
END $$

DELIMITER ;

DELIMITER $$

CREATE trigger payments_update_audit
AFTER UPDATE ON payments FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'orderid', old.orderid, 'docentry', old.docentry, 'docnum', old.docnum, 'canceled', old.canceled, 'paymenttype', old.paymenttype, 'paymentdetails', old.paymentdetails, 'amount', old.amount),
    JSON_OBJECT('id', new.id, 'orderid', new.orderid, 'docentry', new.docentry, 'docnum', new.docnum, 'canceled', new.canceled, 'paymenttype', new.paymenttype, 'paymentdetails', new.paymentdetails, 'amount', new.amount), 
    'UPDATE', null);
END $$

DELIMITER ;

DELIMITER $$

CREATE trigger payments_delete_audit
AFTER DELETE ON payments FOR EACH ROW
BEGIN
	INSERT INTO audits (`old_row_data`, `new_row_data`, `dml_type`, `dml_created_by`) 
    VALUES
    (JSON_OBJECT('id', old.id, 'orderid', old.orderid, 'docentry', old.docentry, 'docnum', old.docnum, 'canceled', old.canceled, 'paymenttype', old.paymenttype, 'paymentdetails', old.paymentdetails, 'amount', old.amount),
	null, 'DELETE', null);
END $$

DELIMITER ;

DROP procedure IF EXISTS `searchDB`;

DELIMITER $$
USE `retail`$$
CREATE PROCEDURE `searchDB` (tableSchema nvarchar(50), searchTerm nvarchar(50))
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

END$$

DELIMITER ;