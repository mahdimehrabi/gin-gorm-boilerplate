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

see docker logs `docker-compose logs -f` and wait for your application to be ready and check your application work by sending get request to `localhost:port/api/ping` (default port is 8000)<br />

### swagger
install [gin swagger](https://github.com/swaggo/gin-swagger) .<br/>
generate swagger docs `swag init` . <br/>
see swagger documention for get information about already implemented rest api services
`localhost:port/swagger/index.html` .

## Features
#### env
we have two env file for this application. </br>
`.env`=> for production and development environment  </br>
`.env.test`=>for test </br>
don't change `Envrionment` key of .env.test it must be `test` </br>
we have an infrastracture for env in `core/infrastracture/env.go` that is responsible for load environment variable in to struct for more clean and easier access around application .
#### dependency injection
we used [uber fx](https://github.com/uber-go/fx) for dependency injection to have more clean application I sugget you read whole documention of fx package
#### docker compose implemention for development and production 
for production use set `Environment=production` and use `docker-compose -f docker-compose.prod.yml up -d ` to run app
for development set `Environment=development` and use `docker-compose up -d ` to run app
#### delve debugger + source watcher for reload   
we have a powerful source watcher for reload server and delve debuger configuration <br />
`docker/dev/web.sh` handle watching and running delve server you can use configuration in `.vscode` directory for configure your vscode to connect to debugger for debugging in normal and  even debugging your tests
#### advanced jwt authentication + middleware
we have an app called authApp that is responsible for jwt authentication and related stuff like sending forgot password email, change password,password strength checker and etc. you can see its rest api services in 
#### saving device name ip and city on login
on users login ,client(frontend) must send login device name and we store that + user IP on login 
and this app have some routes to let user see and manage his logged in devices + information (like IP,country,city) and remove them one by one or all of them (like social media apps)
#### registeration by verifying email
we have complete registeration service for sending email , resending email , verify token and ... that implemented in authentication app
#### swagger documention generator
swagger documention of apps are stored in `core/models/swagger` you can edit them or add new one . 
I suggest you read [swagger gin documention]{https://github.com/swaggo/gin-swagger} for more information.
#### easy custom validator implemention
we have advanced validation system that integerated with gin and [go-playground validation]{https://github.com/go-playground/validator},
for example we have a custom validator that is responsible for checking the of a field is unique in database table, you can use it in your request model like below:
```
type CreateUserRequestAdmin struct {
	IsAdmin        bool   `json:"isAdmin"`
	Email          string `json:"email" binding:"required,uniqueDB=users&email"` //users=>table name , email => column name
}
``` 
you can implement your custom validator by creating your validator in `core/validators/` and introduing them in depndency injection by adding them to `core/validators/validators.go` file in you can understand how to create them by reading the example validators that are exist in validators directory and for more information you can read [gin]{https://github.com/gin-gonic/gin} and [go-playground validation]{https://github.com/go-playground/validator} documents.
#### logger + sentry 
we have logger infrastructre that is responsible for storing logs in sentry I will encourage you to create a [sentry]{https://sentry.io/} account and create a project and add your dsn to .env file `SentryDSN=yourDsn` 
logger is only responsible to capure 500 errors that handled by you in controllers, we configured sentry to capture unwanted panics in `core/bootstrap.go` file.
#### migration versioning 
migrations are store in `core/migrations` and we use this [package]{https://github.com/golang-migrate/migrate} for our migration and we have some commands to generate and run them in Makefile 
#### easy response 
look at `core/response/responses.go` file you can use these functions to have a clean responses for success and error rest response
### middlewares 
you can create your middleware in app in apps directory but if you think the middleware you want to develop is a generic middleware I suggest you put your middleware in generic app <br />
after creating your middleware don't forget to add it do dependency injection , you can add it by editing `core/middlewares` file.
#### database transaction middleware
this middleware is responsible for database transaction you add it to your route and if the resonse of the route is not successful it roolback the transaction.
#### pagination
check out `apps/userApp/repositories/user.go` and `getAllUsers` method
#### cmd
you can create your custom cmd commands by editing `Makefile` and adding your functions to `cmd` package for example we have a command for create-app
#### image resizer 
implemented in `UploadProfilePicture` method of `apps/userApp/controllers/profile.go` 
#### Testing 
an advanced integeration test implemented using [testify]{https://github.com/stretchr/testify} package 
you can use these command for running tests

`make test-all` : run all tests of application 

`make test-all-debugger`: run all tests of application with delve debugger server so you can connect your editor/IDE to it for debugging

`make kill-test-debugger`: kill opened test delve debugger server 

