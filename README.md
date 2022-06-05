# Gin Gorm Boilerplate with clean architecture

A clean archtiecture and well documented gin boilerplate with lot of good features
ready for you to use and improve your speed of development 

## requirements
make sure you installed [docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/)

## Installation 
clone repository in your system
```git clone https://github.com/mahdimehrabi/gin-gorm-boilerplate.git```

create env files base on example 
```
cd gin-gorm-boilerplate
cp env.example .env
cp env.test.example .env.test
```
if you want to run on localhost set `Environment=development` and if you want to run in production model set `Environment=production` in .env file

create docker volume for database data `docker volume create psql_data`

build and run docker-compose  
for development:`docker-compose up -d` <br />
for production: `docker-compose -f docker-compose.prod.yml up -d` 

run migrations 
```
make migrate-up
```

### swagger
see docker logs `docker-compose logs -f` and wait for your application to be ready and check your application work by sending get request to `localhost:port/api/ping` (default port is 8000)<br />

install [gin swagger](https://github.com/swaggo/gin-swagger) 
generate swagger docs `swag init`
see swagger documention for already implemented rest api services
`localhost:port/swagger/index.html`

## Features
#### env
we have two env file for this application
`.env`=> for production and development environment 
`.env.test`=>for test 
don't change `Envrionment` key of .env.test it must be `test`
and we have an infrastracture for env in `core/infrastracture/env.go` that is responsible for load environment variable in to struct for more clean and easier access around application 
#### dependency injection
we used [uber fx](https://github.com/uber-go/fx) for dependency injection to have more clean application I sugget you read whole documention of fx package
#### docker compose implemention for development and production 
for production use set `Environment=production` and use `docker-compose -f docker-compose.prod.yml up -d ` to run app
for development set `Environment=development` and use `docker-compose up -d ` to run app
#### delve debugger + source watcher for reload   
we have a powerful source watcher for reload server and delve debuger configuration <br />
`docker/dev/web.sh` handle watching and running delve server you can use configuration in `.vscode` directory for configure your vscode to connect to debugger for debugging in normal and  even debugging your tests
#### advanced jwt authentication + middleware
we have an called authApp that is responsible for jwt authentication you can see its rest api services in [swagger](#place-1)
#### -saving device name ip and city on login
#### -devices
#### registeration by verifying email
#### forgot password
#### change password route 
#### password strength validation
#### swagger documention generator
#### unique field validation + easy custom validator implemention
#### logger
#### migration versioning 
#### testify implemention
#### dynamic validation errors
#### database transaction middleware
#### pagination
#### image resizer 
#### advanced rest api responses
#### cmd
#### image resizer 