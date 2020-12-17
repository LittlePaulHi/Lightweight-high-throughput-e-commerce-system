#!/bin/bash

stages=( "getallproducts_test" "getAllOrderByAccountID_test" "getAllOrderItemsByOrderID_test" "addCart_test" "getAllCartsByAccountID_test" "editCart_test" "PurchaseFromCarts_test" "integrated_test" )
tag="async"

for directory in "${stages[@]}"
do
    for ((vus=500; vus<=3000;vus+=500))
    do
        echo "Restart apiserver"
        ssh apiserver "TAG=${tag} sh ./autorestart.sh"
        sleep 0.5

        echo "Restart kafka"
        ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose down;docker-compose rm -vfs;docker-compose -f ./docker-compose-single-broker.yml up -d"
        sleep 0.5

        echo "Restart redis"
        ssh redis "cd ~/redis/;docker-compose down;docker-compose rm -vfs;docker-compose up -d"
        sleep 0.5

        echo "Restart consumer"
        ssh consumer "TAG=${tag} sh ./autorestart.sh"
        sleep 0.5

        echo "Restore Database"
        mysql --defaults-extra-file=config < ppfinal.sql
        sleep 0.5

        k6 run -e TIMES=${vus} --out influxdb=http://192.168.0.5:8086/${directory} "../${directory}.js"
        sleep 0.5
    done
done

