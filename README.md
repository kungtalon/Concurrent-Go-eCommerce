## Concurrent-Go-eCommerce
An eCommerce platform supports high volume of visits using Iris framework in Go.

### Basic
1. Used Iris, GORM& MySQL for creating fundamental CRUD services of the eCommerce web app.
2. Used RabbitMQ to control the flow, alleviating the load of the database.
3. Used AES for user authentication with cookie and consistent hashing for distributed user authentication.
4. Used gRPC for communicating between services during authentication.

### Structure
- admin: administrative dashboard for sellers.
- common: utility tools, including reading cookie, connecting to database, preprocessing filter and constants.
- datamodels: data models defined with gorm.Model.
- distributed: tools for distributed authentication, including consistent hash and RabbitMQ.
- front: main web app with user interface.
- lightning: main services for distributed user identification and product number validation.
- repositories: create, read, update, delete (CRUD) functions implemented with GORM.
- services: APIs for interacting with database built on repositories functions.

### Notes

**Q: For message broker, why is RabbitMQ chosen over Kafka?**

- Both RabbitMQ and Apache Kafka are popular open-source message brokers. If the size of every message is small and low probability of data loss is required, RabbitMQ is an ideal choice. In this project, RabbitMQ is used for passing real-time order information to consumers where the database gets updated.  Therefore, RabbitMQ is preferred due to its better reliability and faster real-time messaging.

**Q: What is the difference between REST API and gRPC?**

|                        | REST                   | gRPC                                       |
| ---------------------- | ---------------------- | ------------------------------------------ |
| Scenario               | Web App, Microservices | Between Microservices                      |
| Messaging Format       | JSON, XML              | Protocol Buffer                            |
| Protocol               | HTTP 1.1               | HTTP 2                                     |
| Transmission Speed     | Lower                  | Faster                                     |
| Model of Communication | request-response       | request-response, streaming, bidirectional |
| Method Categories      | GET, POST, PUT ...     | Totally Custom                             |

* gRPC is not suitable for creating web user interface because of its bad browser compatibility.



### Useful References

[1] Go-Iris Documentation: https://docs.iris-go.com/iris/mvc/mvc-quickstart

[2] GORM Documentation: https://gorm.io/docs/index.html

[3] RabbitMQ tutorial in Go: https://www.rabbitmq.com/tutorials/tutorial-one-go.html

[4] Consistent hash: https://akshatm.svbtle.com/consistent-hash-rings-theory-and-implementation & https://segmentfault.com/a/1190000021199728

[5] gRPC with Go: https://grpc.io/docs/languages/go/quickstart/

