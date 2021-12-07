FROM golang:1.17-alpine AS GO_BUILD
COPY . /app
WORKDIR /app
RUN go build -o /go/bin/app scaler/main.go
FROM alpine:3.10
RUN apk update && apk add bash && apk add curl
WORKDIR app
COPY --from=GO_BUILD /go/bin/app ./
CMD ./app
