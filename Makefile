.PHONY: build unit-test rest-integration-test mongo-integration-test kafka-integration-test component-test

build:
	@go build -o bin/events-server cmd/events-server/main.go

unit-test:
	@go test -v -short -coverprofile c.out ./...

rest-integration-test:
	@go test -v -c -o test/integration/http/rest/handler_test test/integration/http/rest/*.go
	./test/integration/http/rest/handler_test -test.v

mongo-integration-test:
	@go test -v -c -o test/integration/mongo/mongo_test test/integration/mongo/*.go
	@docker-compose -f test/integration/mongo/docker-compose.yaml build --no-cache --force-rm
	@docker-compose -f test/integration/mongo/docker-compose.yaml up --abort-on-container-exit --exit-code-from mongo_test

kafka-integration-test:
	@go test -v -c -o test/integration/kafka/kafka_test test/integration/kafka/*.go
	@docker-compose -f test/integration/kafka/docker-compose.yaml build --no-cache --force-rm
	@docker-compose -f test/integration/kafka/docker-compose.yaml up --abort-on-container-exit --exit-code-from kafka_test

component-test:
	@go test -v -c -o test/component/component_test test/component/*.go
	@docker-compose -f test/component/docker-compose.yaml build --no-cache --force-rm
	@docker-compose -f test/component/docker-compose.yaml up --abort-on-container-exit --exit-code-from component_test
