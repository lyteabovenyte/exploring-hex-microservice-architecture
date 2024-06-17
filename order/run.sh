#!/bin/zsh

DATA_SOURCE_URL=root:verysecretpass@tcp(127.0.0.1:3000)/order
APPLICATION_PORT=3000
ENV=development

go run cmd/main.go