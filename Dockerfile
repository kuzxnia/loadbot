# build a small image
FROM alpine:3.19.1
COPY lbot /usr/local/bin/lbot

ENTRYPOINT ["./usr/local/bin/lbot"]
