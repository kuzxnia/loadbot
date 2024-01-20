# base image
FROM golang:1.20.5-alpine AS builder
# create appuser.
RUN adduser -D -g '' user
# create workspace
WORKDIR /opt/app/
COPY go.mod go.sum ./
# fetch dependancies
RUN go mod download && \
  go mod verify
# copy the source code as the last step
COPY . .
# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/api ./cmd/main.go


# build a small image
FROM alpine:3.18
LABEL language="golang"
# import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
# copy the static executable
COPY --from=builder --chown=user:1000 /go/bin/api /api
# use a non-root user
USER user
# run app
ENTRYPOINT ["./api"]
