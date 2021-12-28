# prime-wasm-filter

Sample project to builld a istio/envoy extension using WASM for Rust

Based on:

- [Extending envoy with WASM and Rust](https://antweiss.com/blog/extending-envoy-with-wasm-and-rust/)
- [Proxy WASM Rust](https://github.com/otomato-gh/proxy-wasm-rust)

## Getting started

- Generate extension with docker

```shell
docker build . -t prime-wasm-filter
```

- Copy extension to host using

```shell
docker run -v $PWD/release/wasm32-unknown-unknown/:/opt/mount --rm --entrypoint cp prime-wasm-filter /target/wasm32-unknown-unknown/release/primeenvoyfilter.wasm /opt/mount/primeenvoyfilter.wasm 
```

- Run docker compose with istio envoy with the wasm extension using

```shell
docker-compose -f ./release/docker-compose.yaml up --build
```

In order to test execute the following instruction

- OK `curl  -H "token":"32323" 0.0.0.0:18000`
- FAIL `curl  -H "token":"323232" 0.0.0.0:18000`

## Useful commands

- Connect to docker to browser content using sh
```shell
docker run -it --entrypoint sh prime-wasm-filter
```

The generated file is located in `/target/wasm32-unknown-unknown/release` with the name `primeenvoyfilter.wasm`