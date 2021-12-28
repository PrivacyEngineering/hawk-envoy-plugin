FROM rust:alpine3.14
# flag and build c compiler for wasm32
ENV RUSTFLAGS="-C target-feature=-crt-static"
RUN apk add build-base

# enable target for web assembly
RUN rustup target add wasm32-unknown-unknown

# dependencies
COPY ./Cargo.toml .
COPY ./Cargo.lock .

# source code
COPY src ./src

# build library for target wasm32 optimize for release
RUN cargo build --target=wasm32-unknown-unknown --release
