#!/usr/bin/env bash

echo "BUILDING"
docker build -t gym-tracker .

echo "RUNNING"
docker run \
  --name gym-tracker \
  --rm \
  -d \
  -v ${DOCKER_DATA}/gym-tracker/gcp-auth.json:/root/gcp-auth.json \
  -h puregymtrackerback \
  -e GYMTRACKER_ENV=PRODUCTION \
  -e GYMTRACKER_PORT_BACKEND \
  -e GYMTRACKER_GOOGLE_PROJECT \
  -e GYMTRACKER_INFLUX_URL \
  -e GYMTRACKER_INFLUX_USER \
  -e GYMTRACKER_INFLUX_PASS \
  -e GYMTRACKER_INFLUX_DATABASE \
  -e GYMTRACKER_INFLUX_RETENTION \
  gym-tracker

echo "LOGGING"
docker logs -f gym-tracker
