FROM tinygo/tinygo:0.22.0
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux tinygo build -o actixgo-filter.wasm -scheduler=none -target=wasi ./pkg