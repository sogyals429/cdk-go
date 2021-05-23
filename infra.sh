#!/usr/bin/env bash

(
  rm main
  cd funcs || { echo "Failure"; exit 1; }
  GOOS=linux CGO_ENABLED=0 go build main.go
  zip function.zip main
)

