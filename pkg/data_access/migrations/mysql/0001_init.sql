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

CREATE TABLE IF NOT exists user (
    id int AUTO_INCREMENT,
    user_id varchar(10) UNIQUE,
    password varchar(100),
    name varchar(50),
    age int,
    position varchar(50),
    email varchar(50),
    is_active boolean,
    role varchar(50),
    created_date timestamp,
    PRIMARY KEY (id)
);

INSERT INTO items (id,`name`,price,`type`, img) VALUES
	 (1,'trà đào cam xả',50,1, "1"),
	 (2,'cà phê đen',35,1, "2");

INSERT INTO customer (id,`name`, age, loyalty_point) VALUES
	 (1,'Nguyễn văn a', 50, 10),
	 (2,'trần văn B', 24, 100),
	 (3,'La thị C', 35, 20),
	 (4,'Phạm văn D', 40, 0);

INSERT INTO user (id, user_id, `name`, `password`, age, position, email, is_active, `role`, created_date) VALUES
	 (1, "admin", 'Nguyễn văn a', "$2a$10$1joc.H1g998T5NL2/6as9ugfNmWcx4mzkhPL8hIguIPyye6RSlWP2", 50, "Nhân viên quầy", "anv@gmail.com", 1, "admin", "2014-01-06 18:36:00"),
	 (2, "counter", 'Nguyễn văn B', "$2a$10$.UMihNYTxPbMtlwuTwpc2uxDHDQOzzI6I8E9qq2/fxRyRQFooxX3W", 24, "Nhân viên quầy", "bnv@gmail.com", 1, "counter", "2014-01-06 18:36:00"),
	 (3, "kitchen", 'Nguyễn văn nam', "$2a$10$rt31WFZ2ZJ2p32c09bpxRedefnfY8cEs3d73NLEqtTCEmmeYdUnLK", 28, "bếp trưởng", "namnv@gmail.com", 1, "kitchen", "2014-01-06 18:36:00"),
	 (4, "namnv", 'trần văn nam', "$2a$10$SB/DiVCUMNG/ZYIJkjVEw.ey9EP.VZl7tIfPGDqk5DFISo2yZTdFO", 28, "quản lý cao cấp", "namnv1@gmail.com", 1, "admin", "2014-01-06 18:36:00");

COMMIT;

-- +migrate Down

