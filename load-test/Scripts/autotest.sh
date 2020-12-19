#!/bin/bash

stages=( "getallproducts_test" "getAllOrderByAccountID_test" "getAllOrderItemsByOrderID_test" "addCart_test" "getAllCartsByAccountID_test" "editCart_test" "PurchaseFromCarts_test" "integrated_test" )
tag="async"


for directory in "${stages[@]}"
do
    for ((vus=10; vus<=10;vus+=500))
    do
        echo "Stop kafka"
        ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose stop;docker-compose rm -vfs"
        sleep 0.5

        echo "Stop redis"
        ssh redis "cd ~/redis/;docker-compose stop;docker-compose rm -vfs"

        echo "Stop consumer"
        ssh consumer "TAG=${tag} sh ./stop.sh"

        echo "Stop apiserver"
        ssh apiserver "TAG=${tag} sh ./stop.sh"
        
        echo "Restore Database"
        mysql --defaults-extra-file=config < ppfinal.sql

        echo "Start kafka"
        ssh ppf204@kafka "cd ~/kafka-docker/;docker-compose -f ./docker-compose-single-broker.yml up -d --no-recreate"
        sleep 0.5

        sleep 30s

        echo "Start redis"
        ssh redis "cd ~/redis/;docker-compose up -d"

        echo "Start consumer"
        ssh consumer "TAG=${tag} sh ./start.sh"

        echo "Start apiserver"
        ssh apiserver "TAG=${tag} sh ./start.sh"
        
        sleep 1m

        ./k6 run -e TIMES=${vus} --out influxdb=http://192.168.0.5:8086/${directory} "../${directory}.js"
    done
done


#        sleep 0.5
