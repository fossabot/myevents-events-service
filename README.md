# [MyEvents](https://github.com/danielpacak/myevents) :: Events Service

[![Build Status](https://travis-ci.org/danielpacak/myevents-events-service.svg?branch=master)](https://travis-ci.org/danielpacak/myevents-events-service)
[![codecov](https://codecov.io/gh/danielpacak/myevents-events-service/branch/master/graph/badge.svg)](https://codecov.io/gh/danielpacak/myevents-events-service)

The events service handles the events, their locations, and changes that happen to them.
It's part of the [MyEvents](https://github.com/danielpacak/myevents) application.

## Configuration

| Name                      | Default Value       | Description            |
|---------------------------|---------------------|------------------------|
| MONGODB_CONNECTION_URL    | mongodb://127.0.0.1 | MongoDB connection URL |
| MONGODB_DATABASE_NAME     | myevents            | MongoDB database name  |

## Building and running

1. Download the source code:
   ```
   $ go get https://github.com/danielpacak/myevents-event-service.git
   $ cd $GOPATH/src/github.com/danielpacak/myevents-event-service
   ```
2. Start [MongoDB](https://www.mongodb.com) in Docker container:
   ```
   $ docker run -d --rm -p 27017:27017 --name events-db mongo:latest
   ```
3. Choose a message broker, either [RabbitMQ](https://www.rabbitmq.com/) or
   [Apache Kafka](https://kafka.apache.org/).
   3. RabbitMQ
      3. Start RabbitMQ in Docker container:
         ```
         $ docker run -d --rm -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3-management
         ```
         After starting the container, you will be able to open an AMQP connection to
         `amqp://localhost:5672` and open the management UI in your browser at
         [http://localhost:15672](http://localhost:15672). The default administrator username
         and password are `guest` and `guest`.
   3. Apache Kafka
      4. Start [Apache Zookeeper](https://zookeeper.apache.org/) in Docker container:
         ```
         $ docker run -d --rm --name zookeeper --network host \
           -e ZOOKEEPER_CLIENT_PORT=2181 \
           -e ZOOKEEPER_TICK_TIME=2000 \
           -e ZOOKEEPER_LOG4J_ROOT_LOGLEVEL=ERROR \
           confluentinc/cp-zookeeper:5.1.2
         ```
      5. Start Kafka broker in Docker container:
         ```
         $ docker run -d --rm --name kafka --network host \
           -e KAFKA_ZOOKEEPER_CONNECT="localhost:2181" \
           -e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://localhost:9092" \
           confluentinc/cp-kafka:5.1.2
         ```

### Building and running locally

```
$ go build
$ ./myevents-event-service
```

### Building and running with Docker

```
$ GOOS=linux go build
$ docker image build -t danielpacak/myevents-events-service .
```

```
$ docker container run -d --name events \
     -e AMQP_BROKER_URL=amqp://guest:guest@localhost:5672/ \
     -e MONGO_URL=mongodb://localhost:27017/events \
     -p 8181:8181 \
     -p 9100:9100 \
     danielpacak/myevents-events-service:latest
```
