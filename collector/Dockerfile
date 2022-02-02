FROM golang:1.17-buster as builder
WORKDIR /app
COPY . .
RUN go build -o /app/collector
EXPOSE 8080
CMD [ "/app/collector" ]

#FROM gcr.io/distroless/static-debian11
#WORKDIR /
#COPY --from=builder /app/collector /collector
#EXPOSE 8080
#USER nonroot:nonroot
#ENTRYPOINT "/collector"
