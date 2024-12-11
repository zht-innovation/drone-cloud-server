#!/bin/bash
wget -P ~/install/ https://github.com/redis/redis/archive/refs/tags/7.4.1.tar.gz
cd ~/install
tar xvzf redis-7.4.1.tar.gz
cd redis-7.4.1
make
sudo make install
sudo redis-server /etc/redis/redis.conf


# method2: sudo apt install redis-server