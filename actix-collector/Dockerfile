FROM rust:alpine3.14 as builder
# dependencies
COPY ./Cargo.toml .
COPY ./Cargo.lock .

# source code
COPY src ./src

# build for release
RUN cargo build --release

#FROM alpine:3.14
#COPY --from=builder