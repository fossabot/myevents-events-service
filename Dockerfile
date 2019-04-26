FROM scratch

LABEL maintainer="pacak.daniel@gmail.com"
LABEL name="myevents-events-service"

COPY bin/events-server /events-server

ENTRYPOINT ["/events-server"]
