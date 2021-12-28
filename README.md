# rust-envoy-proxy

[![Apache 2.0 License][license-badge]][license-link]

[license-badge]: https://img.shields.io/github/license/proxy-wasm/proxy-wasm-rust-sdk
[license-link]: https://github.com/TUB-CNPE-TB/rust-envoy-proxy/blob/master/LICENSE

Extension for istio envoy to allow trace personal data between rest microservices in kubernetes

## Projects

### [prime-wasm-filter](./prime-wasm-filter)

Sample project to use a prime factor filter to authorize the flow.
It uses rust to generate a wasm32 using docker, and then install the extension in
envoy proxy. Test evverything using docker compose

## Rust programming language documentation

- [Rust - The Book](https://doc.rust-lang.org/book/)
- [Rust by example](https://doc.rust-lang.org/stable/rust-by-example/)
- [REST api framework for rust](https://actix.rs/docs/)

Blogs

- [Extending envoy with wasm and rust](https://antweiss.com/blog/extending-envoy-with-wasm-and-rust/)
- [Istio wasm plugin configuration](https://istio.io/latest/docs/reference/config/proxy_extensions/wasm-plugin/)
- [Istio: Distributing WebAssembly Modules](https://istio.io/latest/docs/ops/configuration/extensibility/wasm-module-distribution/)

## Videos

- [WebAssembly Extension and envoy architecture](https://www.youtube.com/watch?v=XdWmm_mtVXI)
- [WebAssembly extension for proxies](https://www.youtube.com/watch?v=OIUPf8m7CGA)
