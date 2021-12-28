docker build . -t prime-wasm-filter
docker run -v $PWD/release/wasm32-unknown-unknown/:/opt/mount --rm --entrypoint cp prime-wasm-filter /target/wasm32-unknown-unknown/release/primeenvoyfilter.wasm /opt/mount/primeenvoyfilter.wasm
sha256sum release/wasm32-unknown-unknown/primeenvoyfilter.wasm
docker-compose -f ./release/docker-compose.yaml up --build
docker-compose -f ./release/docker-compose.yaml rm