#!/bin/bash
docker run -d -p 3080:9100 api-service:${TAG} --restart-policy=restart