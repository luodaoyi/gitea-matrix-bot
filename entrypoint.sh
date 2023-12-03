#!/bin/sh
printf "nameserver 127.0.0.11\nnameserver 8.8.4.4\nnameserver 223.5.5.5\n" > /etc/resolv.conf

# 检查sqlite是否存在
if [ ! -f /data/tokens.db ]; then
    # 如果不存在则创建
    ./app initdb 
fi
# 启动服务
./app -v 
