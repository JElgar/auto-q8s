FROM golang:1.17-alpine AS GO_BUILD
COPY . /app
WORKDIR /app
RUN go build -o /go/bin/server consumer/consumer.go
FROM alpine:3.10
WORKDIR app
COPY --from=GO_BUILD /go/bin/server ./
EXPOSE 10000
CMD ./server
