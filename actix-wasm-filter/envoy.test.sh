docker build . -t actix-wasm-filter
docker run -v $PWD/release/wasm32-unknown-unknown/:/opt/mount --rm --entrypoint cp actix-wasm-filter /target/wasm32-unknown-unknown/release/actixenvoyfilter.wasm /opt/mount/actixenvoyfilter.wasm
sha256sum release/wasm32-unknown-unknown/actixenvoyfilter.wasm
docker-compose -f ./release/docker-compose.yaml up --build
docker-compose -f ./release/docker-compose.yaml rm