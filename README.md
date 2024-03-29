# [MyEvents](https://github.com/danielpacak/myevents) :: Events Service

[![Build Status](https://travis-ci.org/danielpacak/myevents-events-service.svg?branch=master)](https://travis-ci.org/danielpacak/myevents-events-service)
[![Maintainability](https://api.codeclimate.com/v1/badges/dcb05fb58271dc0f5006/maintainability)](https://codeclimate.com/github/danielpacak/myevents-events-service/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/dcb05fb58271dc0f5006/test_coverage)](https://codeclimate.com/github/danielpacak/myevents-events-service/test_coverage)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdanielpacak%2Fmyevents-events-service.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdanielpacak%2Fmyevents-events-service?ref=badge_shield)

The events service handles the events, their locations, and changes that happen to them.
It's part of the [MyEvents](https://github.com/danielpacak/myevents) application.

## Architecture

![Architecture](docs/architecture.png)

## Configuration

| Name                                | Default Value         | Description            |
|-------------------------------------|-----------------------|------------------------|
| `EVENTS_SVC_MONGODB_CONNECTION_URL` | `mongodb://127.0.0.1` | MongoDB connection URL |
| `EVENTS_SVC_MONGODB_DATABASE_NAME`  | `events-svc`          | MongoDB database name  |
| `EVENTS_SVC_KAFKA_BROKERS`          | `localhost:9092`      | ? |
| `EVENTS_SVC_OUTBOUND_TOPIC`         | `myevents.events`     | ? |
| `EVENTS_SVC_REST_API_TCP_ADDRESS`   | `:8181`               | ? |
| `EVENTS_SVC_METRICS_TCP_ADDRESS`    | `:9100`               | ? |

## Building and running

1. Download the source code:
   ```
   $ go get https://github.com/danielpacak/myevents-event-service.git
   $ cd $GOPATH/src/github.com/danielpacak/myevents-event-service
   ```
2. Build executable:
   ```
   $ GOOS=linux go build
   ```
3. Start the backing services:
   ```
   $ docker-compose up -d
   ```
4. Create a sample event with curl:
   ```
   curl -0 -v http://localhost:8181/events \
   -H 'Content-Type: text/json; charset=utf-8' \
   -d @- << EOF
   {
       "name": "Dzień Niepodległości Armenii",
       "duration": 5,
       "location": {
           "Name": "PKiN",
           "Address": "Plac Defilad 1, 00-901 Warszawa",
           "Country": "Poland",
           "OpenTime": 30,
           "EndTime": 40,
           "Halls": [{
               "name": "Main hall",
               "location": "1st floor",
               "capacity": 150
           }]
       }
   }
   EOF
   ```
5. Navigate to Mongo Express at [http://localhost:8081](http://localhost:8081) and make sure
   that the corresponding document was created in events collection.
6. Navigate to Control Center at [http://localhost:9021](http://localhost:9021) and make sure
   that the corresponding message was published to the outbound topic.

## References

1. [Writing a Go client for your RESTful API](https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c)
2. [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
3. [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)
4. [Integration Test With Database in Golang](https://hackernoon.com/integration-test-with-database-in-golang-355dc123fdc9)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdanielpacak%2Fmyevents-events-service.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdanielpacak%2Fmyevents-events-service?ref=badge_large)