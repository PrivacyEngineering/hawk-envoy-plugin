# actix-collector

Sample project in rust to run a rest api and receive any request from any server

## Getting started

1. Build the project docker image and publish to docker hub

```shell
docker build . -t actix-collector -t jmgoyesc/actix-collector
docker push jmgoyesc/actix-collector
```

2. Run the project

```shell
docker run -p 8080:8080 --name run-actix-colletor jmgoyesc/actix-collector
```

3. Create kubernetes deployment and service

```shell
kaf deploy/
```
