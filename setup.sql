CREATE TABLE user (
    id          INT(10)  AUTO_INCREMENT PRIMARY KEY NOT NULL,
    user_name        VARCHAR(30) NOT NULL,

);

CREATE TABLE vehicle (
    id          INT(10)  AUTO_INCREMENT PRIMARY KEY NOT NULL,
    name        VARCHAR(30) NOT NULL,
    type        VARCHAR(30)
    description VARCHAR(255),
    images      VARCHAR(255),
    daily_price INT(10) ,       
    status      boolean

);



CREATE TABLE booking (
     id           INT(10)  AUTO_INCREMENT PRIMARY KEY NOT NULL,
     user_id      INT(10), CONSTRAINT FOREIGN KEY 	(user_id) REFERENCES user(id)
     vehicle_id   INT(10), CONSTRAINT FOREIGN KEY 	(vehicle_id) REFERENCES vehicle(id)
     totalconst   INT(10),
     start_date   DATE,
     end_date     DATE
);

CREATE TABLE wallet (
    id           INT(10)  AUTO_INCREMENT PRIMARY KEY NOT NULL,
    user_id      INT(10)  CONSTRAINT FOREIGN KEY 	(user_id) REFERENCES user(id)
    value        INT(10)
)

