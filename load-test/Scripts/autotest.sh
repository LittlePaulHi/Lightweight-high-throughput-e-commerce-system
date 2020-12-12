#!/bin/bash

stages=( "getallproducts_test" "getAllOrderByAccountID_test" "getAllOrderItemsByOrderID_test" "addCart_test" "getAllCartsByAccountID_test" "editCart_test" "PurchaseFromCarts_test" "integrated_test" )

ssh apiserver "sh ./autorestart.sh"
ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
ssh redis "cd ~/redis/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
ssh consumer "sh ./autorestart.sh"

for directory in "${stages[@]}"
do
    for filename in $(find ${directory}/*.js | sort -z)
    do
        db=$(echo ${filename} | sed -e "s/\/.*//g" | sed -e "s/.js//g")
        k6 run --out influxdb=http://192.168.0.5:8086/${db} "${filename}"
        ssh apiserver "sh ./autorestart.sh"
        ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
        ssh redis "cd ~/redis/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
        ssh consumer "sh ./autorestart.sh"
    done
done