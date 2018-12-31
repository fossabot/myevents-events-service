# myevents-event-service

[![Build Status](https://travis-ci.org/danielpacak/myevents-event-service.svg?branch=master)](https://travis-ci.org/danielpacak/myevents-event-service)

```
$ go get https://github.com/danielpacak/myevents-event-service.git
$ cd $GOPATH/src/github.com/danielpacak/myevents-event-service
$ go build
$ ./myevents-event-service
```

```
$ docker run -d --rm -p 27017:27017 --name mongodb mongo:latest
```

```
$ docker run -d --rm -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3-management
```

After starting the container, you will be able to open an AMQP connection to
amqp://localhost:5672 and open the management UI in your browser at
[http://localhost:15672](http://localhost:15672).

## TLS

```
$ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365
```

```
$ openssl rsa -in key.pem -out key.unencrypted.pem -passin pass:YOUR_PASSWORD
```
