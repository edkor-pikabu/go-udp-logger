CREATE DATABASE server;

use server;

CREATE TABLE server.logs(
	id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	record_name varchar(150) NOT NULL,
	group_name varchar(150) NOT NULL,
	data text NOT NULL
);