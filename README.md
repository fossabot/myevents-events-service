# [MyEvents](https://github.com/danielpacak/myevents) :: Events Service

[![Build Status](https://travis-ci.org/danielpacak/myevents-events-service.svg?branch=master)](https://travis-ci.org/danielpacak/myevents-events-service)
[![codecov](https://codecov.io/gh/danielpacak/myevents-events-service/branch/master/graph/badge.svg)](https://codecov.io/gh/danielpacak/myevents-events-service)

The events service handles the events, their locations, and changes that happen to them.
It's part of the [MyEvents](https://github.com/danielpacak/myevents) application.

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
3. Start [RabbitMQ](https://www.rabbitmq.com/) in Docker container:
   ```
   $ docker run -d --rm -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3-management
   ```
   After starting the container, you will be able to open an AMQP connection to
   `amqp://localhost:5672` and open the management UI in your browser at
   [http://localhost:15672](http://localhost:15672). The default administrator username
   and password are `guest` and `guest`.

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
