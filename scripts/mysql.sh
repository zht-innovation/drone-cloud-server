#!/bin/bash

sudo apt update
sudo apt install mysql-server -y
sudo mysql -u root -e "
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'asd1234567@';
FLUSH PRIVILEGES;
CREATE DATABASE zht_cloud_db;
"