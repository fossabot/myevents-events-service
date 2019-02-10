# myevents-event-service

[![Build Status](https://travis-ci.org/danielpacak/myevents-events-service.svg?branch=master)](https://travis-ci.org/danielpacak/myevents-events-service)

## Building and running locally

```
$ go get https://github.com/danielpacak/myevents-event-service.git
$ cd $GOPATH/src/github.com/danielpacak/myevents-event-service
$ go build
$ ./myevents-event-service
```

## Building and running with Docker

```
$ GOOS=linux go build
$ docker image build -t danielpacak/myevents-events-service .
```

```
$ docker run -d --rm -p 27017:27017 --name events-db mongo:latest
```

```
$ docker run -d --rm -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3-management
```

```
$ docker container run -d --name events \
     -e AMQP_BROKER_URL=amqp://guest:guest@localhost:5672/ \
     -e MONGO_URL=mongodb://localhost:27017/events \
     -p 8181:8181 \
     danielpacak/myevents-events-service:latest
```

After starting the container, you will be able to open an AMQP connection to
amqp://localhost:5672 and open the management UI in your browser at
[http://localhost:15672](http://localhost:15672).
