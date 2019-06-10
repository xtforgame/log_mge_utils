FROM alpine:3.9

COPY logwatcher /usr/bin
COPY web /usr/logwatcher/web/
RUN chmod +x /usr/bin/logwatcher

RUN apk add --update --no-cache ca-certificates fuse

WORKDIR /usr/logwatcher

# ENTRYPOINT ["/usr/bin/logwatcher"]
CMD ["/usr/bin/logwatcher"]
