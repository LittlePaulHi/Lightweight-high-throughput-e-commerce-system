#!/bin/bash

stages=( "getallproducts_test" "getAllOrderByAccountID_test" "getAllOrderItemsByOrderID_test" "addCart_test" "getAllCartsByAccountID_test" "editCart_test" "PurchaseFromCarts_test" "integrated_test" )
tags=( "async" "sync" )

for ((partitions=64;partitions<=512;partitions+=64))
do
    for directory in "${stages[@]}"
    do
        for tag in "${tags[@]}"
        do
            for ((vus=200; vus<=1000;vus+=200))
            do
                echo ${partitions} ", " ${directory} ", " ${tag} ", " `date` ", start" >> result.txt
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
                ssh ppf204@kafka "cd ~/kafka-docker/; echo \"syncBuyEventTopic:${tag}:1\" > .env_file ;docker-compose -f ./docker-compose-single-broker.yml up -d --no-recreate"
                sleep 0.5
        
                sleep 30s
        
                echo "Start redis"
                ssh redis "cd ~/redis/;docker-compose up -d"
        
                echo "Start consumer"
                ssh consumer "TAG=${tag} sh ./start.sh"
        
                echo "Start apiserver"
                ssh apiserver "TAG=${tag} sh ./start.sh"
                
                sleep 2m
        
                ./k6 run -e TIMES=${vus} --out influxdb=http://192.168.0.5:8086/${directory} "../${directory}.js"
                echo ${partitions} ", " ${directory} ", " ${tag} ", " `date` ", ended" >> result.txt
            done
        done
    done
done

#        sleep 0.5
