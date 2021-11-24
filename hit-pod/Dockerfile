FROM alpine:latest

RUN apk update && apk add curl
ADD hit.sh /usr/local/bin/hit.sh

ENTRYPOINT ["/bin/sh", "-c", "/usr/local/bin/hit.sh"]