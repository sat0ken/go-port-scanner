#!/bin/bash

CGO_ENABLED=1 GOOS=linux go build $1
