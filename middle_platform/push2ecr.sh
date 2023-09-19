#!/bin/bash

TAG=$(git log -1 --format="%cd-%h" --date=short)
NAME=data_app
REPO=032664146980.dkr.ecr.us-east-1.amazonaws.com
# git checkout develop
# git stash
# git pull --rebase

# login first
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $REPO
docker build -t $NAME:$TAG -f ./Dockerfile .

docker tag $NAME:$TAG $REPO/$NAME:$TAG
docker tag $NAME:$TAG $REPO/$NAME:latest
docker push $REPO/$NAME:$TAG
docker push $REPO/$NAME:latest
