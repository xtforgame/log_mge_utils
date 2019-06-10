FROM alpine:3.9

COPY build/alpine3.9/logwatcher /usr/bin
COPY dist/alpine3.9/logwatcher/web /usr/logwatcher/web/
RUN chmod +x /usr/bin/logwatcher

RUN apk add --update --no-cache ca-certificates fuse

# ENTRYPOINT ["/usr/bin/logwatcher"]
