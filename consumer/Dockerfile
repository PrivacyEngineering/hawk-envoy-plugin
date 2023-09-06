FROM golang:1.17-buster as builder
WORKDIR /app
COPY . .
RUN go build -o /app/consumer
CMD [ "/app/consumer" ]
