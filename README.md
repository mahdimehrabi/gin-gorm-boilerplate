# Gin Gorm Boilerplate with clean architecture

## Introduction
A clean archtiecture and well documented gin boilerplate with lot of good features
ready for you to use and improve your speed of development 

## requirements
make sure you installed [docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/)

## Installation 
clone repository in your system
```git clone https://github.com/mahdimehrabi/gin-gorm-boilerplate.git```

create env files base on example 
```
cp env.example .env
cp env.test.example .env.test
```
if you want to run on localhost set `Environment=development` and if you want to run in production model set `Environment=production` in .env file

build and run docker-compose  
development: `docker-compose up -d`
production: `docker-compose -f docker-compose.prod.yml up -d` 

run migrations 
```
make migrate-up
```

see docker logs `docker-compose logs -f` and wait for your application to be ready and check your application work by sending get request to localhost:port/api/ping (default port is 8000)



## Features
clean architecture
env
dependency injection
docker compose implemention for development and production + cache
delve debugger + source watcher for reload   
advanced jwt authentication + middleware
-saving device name ip and city on login
-devices
registeration by verifying email
forgot password
change password route 
password strength validation
swagger documention generator
unique field validation + easy custom validator implemention
logger
migration versioning 
testify implemention
dynamic validation errors
database transaction middleware
pagination
image resizer 
advanced rest api responses
cmd
image resizer 