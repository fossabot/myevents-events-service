FROM scratch

COPY myevents-events-service /events-service

ENTRYPOINT ["/events-service"]
