#!/bin/bash

stages=( "getallproducts_test" "getAllOrderByAccountID_test" "getAllOrderItemsByOrderID_test" "addCart_test" "getAllCartsByAccountID_test" "editCart_test" "PurchaseFromCarts_test" "integrated_test" )

for directory in "${stages[@]}"
do
    for filename in $(find ${directory}/*.js | sort -z)
    do  
        k6 run "${filename}"
        ssh apiserver "sh ./autorestart.sh"
    done
done
