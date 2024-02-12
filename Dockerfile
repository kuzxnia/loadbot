# build a small image
FROM alpine:3.19.1
COPY loadbot /usr/local/bin/loadbot

ENTRYPOINT ["./usr/local/bin/loadbot"]
