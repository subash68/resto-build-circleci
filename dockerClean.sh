#!/bin/bash

### WARNING: THIS WILL REMOVE EVERYTHING FROM DEV MACHINE ###

docker system prune -a

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

docker volume prune

docker rmi $(docker images -a -q)