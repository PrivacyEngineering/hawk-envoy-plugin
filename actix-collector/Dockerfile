FROM rust:1.57.0-buster as builder

# dependencies
COPY ./Cargo.toml .
COPY ./Cargo.lock .

# source code
COPY src ./src

# build for release
RUN cargo build --release

FROM ubuntu:22.04
COPY --from=builder /target/release/actix-collector ./actix-collector
EXPOSE 8080
ENTRYPOINT ./actix-collector