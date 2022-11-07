FROM golang:1.18-alpine

RUN mkdir /app
WORKDIR /app

ADD .. ./

RUN go build -o config-controller ./cmd/config-controller/

EXPOSE 8080

CMD ["./config-controller"]
