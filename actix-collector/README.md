# actix-collector

Sample project in rust to run a rest api and receive any request from any server

## Getting started

1. Build the project

```shell
docker build . -t actix-collector
```

2. Run the project

```shell
docker run -p 8080:8080 --name run-actix-colletor actix-collector
```

