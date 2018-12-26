# go-rest-api-seed

[![Build Status](https://travis-ci.org/danielpacak/go-rest-api-seed.svg?branch=master)](https://travis-ci.org/danielpacak/go-rest-api-seed)

```
$ go get https://github.com/danielpacak/go-rest-api-seed.git
$ cd $GOPATH/src/github.com/danielpacak/go-rest-api-seed
$ go build gorestapi.go
$ ./gorestapi
```

```
$ docker run -d --rm -p 27017:27017 --name mongodb mongo:latest
```

```
$ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365
```

```
$ openssl rsa -in key.pem -out key.unencrypted.pem -passin pass:YOUR_PASSWORD
```
