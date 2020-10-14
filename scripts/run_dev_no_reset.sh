#!/bin/bash

trap killgroup SIGINT

killgroup(){
  echo killing...
  kill 0
}

go run -tags debug main.go serve-api &

wait