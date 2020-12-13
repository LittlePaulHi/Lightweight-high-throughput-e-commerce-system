#!/bin/bash

stages=( "getallproducts_test" "getAllOrderByAccountID_test" "getAllOrderItemsByOrderID_test" "addCart_test" "getAllCartsByAccountID_test" "editCart_test" "PurchaseFromCarts_test" "integrated_test" )
tag="async"

echo "Restart apiserver"
ssh apiserver "TAG=${tag} sh ./autorestart.sh"
echo "Restart kafka"
ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose down;docker-compose rm -vfs;docker-compose -f ./docker-compose-single-broker.yml up -d"
echo "Restart redis"
ssh redis "cd ~/redis/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
echo "Restart consumer"
ssh consumer "TAG=${tag} sh ./autorestart.sh"
echo "Restore Database"
mysql --defaults-extra-file=config < ppfinal.sql

for directory in "${stages[@]}"
do
    for ((vus=500; vus<=3000;vus+=500))
    do
        k6 run -e TIMES=${vus} --out influxdb=http://192.168.0.5:8086/${directory} "../${directory}.js"
        echo "Restart apiserver"
        ssh apiserver "TAG=${tag} sh ./autorestart.sh"
        echo "Restart kafka"
        ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose down;docker-compose rm -vfs;docker-compose -f ./docker-compose-single-broker.yml up -d"
        echo "Restart redis"
        ssh redis "cd ~/redis/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
        echo "Restart consumer"
        ssh consumer "TAG=${tag} sh ./autorestart.sh"
        echo "Restore Database"
        mysql --defaults-extra-file=config < ppfinal.sql
    done
done
