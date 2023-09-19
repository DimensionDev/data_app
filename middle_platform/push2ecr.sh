#!/bin/bash

NAME=data_app
REPO=032664146980.dkr.ecr.us-east-1.amazonaws.com
# git checkout develop
# git stash
# git pull --rebase

# login first
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $REPO
docker build -t $NAME -f .

docker tag $NAME:latest $REPO/$NAME:latest
docker push $REPO/$NAME:latest
