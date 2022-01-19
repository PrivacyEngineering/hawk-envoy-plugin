FROM tinygo/tinygo:0.21.0
COPY . .
RUN tinygo build -o actixgo-filter.wasm -scheduler=none -target=wasi main.go