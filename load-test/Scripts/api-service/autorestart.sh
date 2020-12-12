#!/bin/bash
ID=$(docker ps | awk 'NR==2{print $1}')
docker stop ${ID}
docker run -d -p 3080:9100 api-service --restart-policy=restart