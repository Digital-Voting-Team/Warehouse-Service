FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/warehouse-service
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/warehouse-service /go/src/warehouse-service


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/warehouse-service /usr/local/bin/warehouse-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["warehouse-service"]
