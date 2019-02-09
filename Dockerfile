FROM debian:jessie

COPY myevents-event-service /events-service
RUN useradd events-service
USER events-service

ENV LISTEN_URL=0.0.0.0:8181
EXPOSE 8181
CMD ["/events-service"]