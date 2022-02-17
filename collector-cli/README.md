# collector cli

## Getting started

```shell
k scale --replicas=0 deploy -n hawk-ns collector-cli-deployment
docker build . -t jmgoyesc/collector-cli-go
docker push jmgoyesc/collector-cli-go
k scale --replicas=1 deploy -n hawk-ns collector-cli-deployment
#kaf deploy
```