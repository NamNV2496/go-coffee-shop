-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS coffee;

CREATE TABLE IF NOT exists customer (
    id int AUTO_INCREMENT,
    name varchar(50),
    age int,
    loyalty_point int,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT exists items (
    id int AUTO_INCREMENT,
    name varchar(50),
    price int,
    type int,
    img varchar(255), 
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT exists orders (
    id int AUTO_INCREMENT,
    customer_id int,
    total_amount int,
    status int,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT exists order_item (
    id int AUTO_INCREMENT,
    order_id int,
    item_id int,
    quantity int,
    price int,
    PRIMARY KEY (id)
);

INSERT INTO items (id,`name`,price,`type`, img) VALUES
	 (1,'thịt kho',50,1, "abcd"),
	 (2,'thịt kho hột vịt',150,1, "abcd"),
	 (3,'cá rán',75,2, "abcd"),
	 (4,'rau luộc',40,3, "abcd");

INSERT INTO customer (id,`name`, age, loyalty_point) VALUES
	 (1,'Nguyễn văn a', 50, 1),
	 (2,'trần văn B', 24, 100),
	 (3,'La thị C', 35, 20),
	 (4,'Phạm văn D', 40, 3);

COMMIT;

-- +migrate Down

