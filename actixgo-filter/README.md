# actixgo-filter

WASM istio/envoy filter based on go.

Based on:

- [WASM Proxy for Go lang](https://github.com/tetratelabs/proxy-wasm-go-sdk)

## Getting started

- Generate extension with docker

```shell
docker build . -t actixgo-filter
```

- Copy extension to host using

```shell
docker run -v $PWD/release/:/opt/mount/ --rm --entrypoint cp actixgo-filter /app/actixgo-filter.wasm /opt/mount/actixgo-filter.wasm 
```

- Generate checksum for installer

```shell
sha256sum release/actixgo-filter.wasm
```

All together

```shell
docker build . -t actixgo-filter
docker run -v $PWD/release/:/opt/mount/ --rm --entrypoint cp actixgo-filter /app/actixgo-filter.wasm /opt/mount/actixgo-filter.wasm 
sed -i '' -e "s/sha256: .*/sha256: $(sha256sum release/actixgo-filter.wasm | head -c 64)/g" release/istio/filter/actixgo.filter.config.yaml
kaf release/istio/filter

```

## Istio example

It is required to use istio gateway for the traffic because the http filter is applied for
gateway. It is possible to apply it for inbound or outbound proxy traffic (envoy) but it should go
through the gateway for the filter to work.

1. Run the Google cloud setup with istio enable
2. Create namespace for httpbin demo project

```shell
kaf release/istio/httbin/httbin.gateway.ns.yaml
```

3. Install httpbin

```shell
kaf https://raw.githubusercontent.com/istio/istio/release-1.12/samples/httpbin/httpbin.yaml -n httpbin-gateway
```

4. Configure istio gateway for httpbin project

```shell
kaf release/istio/httbin/istio.gateway.httpbin.yaml
```

5. Export variables to access istio ingress

```shell
export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
```

6. Test the isitio ingress gateway before install the filter. It should get 200 OK

```shell
curl -v -s -I "http://$INGRESS_HOST:$INGRESS_PORT/headers"
```

7. Install actix filter

```shell
kaf release/istio/filter
```

## Useful commands

- Connect to docker to browser content using sh
```shell
docker run -it --entrypoint sh actixgo-filter
```

The generated file is located in `/` with the name `actixgo-filter.wasm`

- Expose wasm plugin locally using docker nginx
```shell
docker run --name some-nginx -p 9090:80 -v $(pwd)/release/:/usr/share/nginx/html:ro -d nginx
```

- How to check ports in use in linux
```shell
lsof -i -P -n | grep LISTEN
```