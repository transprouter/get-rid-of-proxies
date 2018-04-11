#!/bin/bash

docker container rm --force $(docker container ls --filter name=transprouter_ --format {{.ID}})
docker network rm $(docker network ls --filter name=transprouter_ --format {{.ID}})
