#!/bin/bash
docker run  -d consumer-service:${TAG} --restart-policy=restart