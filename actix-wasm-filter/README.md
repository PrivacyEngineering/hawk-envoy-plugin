# actix-wasm-filter

istio/envoy extension using WASM for Rust to intercept request and response traffic for
services using istio ingress. HTTP (L7) extensions.

Based on:

- [Extending envoy with WASM and Rust](https://antweiss.com/blog/extending-envoy-with-wasm-and-rust/)
- [Proxy WASM Rust](https://github.com/otomato-gh/proxy-wasm-rust)
- [Proxy-Wasm vNEXT ABI (Application Binary Interface) specification](https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT)
- [Envoy Proxy: Life of a request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request#request-flow)
- [Envoy Proxy: Tracing Request ID generation](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/observability/tracing#arch-overview-tracing)
- [Envoy Proxy: x-request-id header](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#config-http-conn-man-headers-x-request-id)
- [Envoy Proxy: Lua filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/lua_filter#config-http-filters-lua-stream-handle-api)

## Getting started

- Generate extension with docker

```shell
docker build . -t actix-wasm-filter
```

- Copy extension to host using

```shell
docker run -v $PWD/release/wasm32-unknown-unknown/:/opt/mount --rm --entrypoint cp actix-wasm-filter /target/wasm32-unknown-unknown/release/actixenvoyfilter.wasm /opt/mount/actixenvoyfilter.wasm 
```

- Generate checksum for installer

```shell
sha256sum release/wasm32-unknown-unknown/actixenvoyfilter.wasm
```

- Replace generated checksum in istio actix.filter.yaml

- Run docker compose with istio envoy with the wasm extension using

```shell
docker-compose -f ./release/docker-compose.yaml up --build -d
```

In order to test execute the following instruction

- OK `curl  -H "x-prime-token":"32323" 0.0.0.0:18000`
- FAIL `curl  -H "x-prime-token":"323232" 0.0.0.0:18000`

Shutdown docker compose

```shell
docker-compose -f ./release/docker-compose.yaml stop
docker-compose -f ./release/docker-compose.yaml rm
```

## Istio example

It is required to use istio gateway for the traffic because the http filter is applied for
gateway. It is possible to apply it for inbound or outbound proxy traffic (envoy) but it should go
through the gateway for the filter to work.

1. Run the Google cloud setup with istio enable
2. Create namespace for httpbin demo project

```shell
kaf release/istio/httbin.gateway.ns.yaml
```

3. Install httpbin

```shell
kaf https://raw.githubusercontent.com/istio/istio/release-1.12/samples/httpbin/httpbin.yaml -n httpbin-gateway
```

4. Configure istio gateway for httpbin project

```shell
kaf release/istio/istio.gateway.httpbin.yaml
```

5. Run release/istio/ files to install the filter. Istio will install the filter in each envoy proxy
6. Export variables to access istio ingress

```shell
export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
```

7. Test the isitio ingress gateway before install the filter. It should get 200 OK

```shell
curl -v -s -I "http://$INGRESS_HOST:$INGRESS_PORT/headers"
```

8. Install prime filter

```shell
kaf release/istio/filter
```

9. Execute the test

HTTP/1.1 403 Forbidden
```shell
curl -H "x-prime-token":"3232" -v -s -I "http://$INGRESS_HOST:$INGRESS_PORT/headers"
```

HTTP/1.1 200 OK
```shell
curl -H "x-prime-token":"32323" -v -s -I "http://$INGRESS_HOST:$INGRESS_PORT/headers"
```

10. Delete prime filter

```shell
k delete -f release/istio/filter
```

## Useful commands

- Connect to docker to browser content using sh
```shell
docker run -it --entrypoint sh actix-wasm-filter
```

The generated file is located in `/target/wasm32-unknown-unknown/release` with the name `actixenvoyfilter.wasm`