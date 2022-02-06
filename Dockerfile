FROM alpine:latest

COPY bin/server/store.sh /usr/bin/store.sh
VOLUME [ "/database" ]
EXPOSE 8080
ENTRYPOINT [ "store.sh" ]