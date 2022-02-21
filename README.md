## Concurrent-Go-eCommerce
An eCommerce platform supports high concurrency using Iris framework in Go.

### Basic:
1. Used Iris, Gorm & MySQL for creating fundamental CRUD services of the eCommerce web app.
2. Used RabbitMQ to control the flow, alleviating the load of the database.
3. Used consistent hash for distributed user authentication.
4. Used gRPC for communicating between services during authentication.

### Structure
admin: administrative dashboard for sellers.
common: utilitary tools, including reading cookie, conneting to database, preprocessing filter and consts.
datamodels: data models defined with gorm.Model.
distributed: tools for distributed authentication, including consistent hash and rabbitmq.
front: main web app with user interface.
lightning: main services for distributed user identification and product number validation.
repositories: create, read, update, delete (CRUD) functions implemented with gorm.
services: APIs for interacting with database built on repositories functions.

### Q & A
To be updated.
