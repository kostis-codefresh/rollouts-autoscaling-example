FROM golang:1.24-alpine3.21 AS build-env

WORKDIR /tmp/workdir

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build 

FROM alpine:3.18

EXPOSE 8080

RUN apk add --no-cache ca-certificates bash

COPY --from=build-env /tmp/workdir/autoscale-example /app/autoscale-example

WORKDIR /app

CMD ["./autoscale-example"]
