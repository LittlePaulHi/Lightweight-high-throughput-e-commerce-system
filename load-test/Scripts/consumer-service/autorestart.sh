#!/bin/bash
ID=$(docker ps | awk 'NR==2{print $1}')
docker stop ${ID}
docker run  -d consumer-service:${TAG} --restart-policy=restart