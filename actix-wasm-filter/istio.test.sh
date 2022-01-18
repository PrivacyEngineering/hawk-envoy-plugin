# create httpbin and apply gateway
kaf release/istio/httbin.gateway.ns.yaml
kaf https://raw.githubusercontent.com/istio/istio/release-1.12/samples/httpbin/httpbin.yaml -n httpbin-gateway
kaf release/istio/istio.gateway.httpbin.yaml

# export istio gateway service host and port
export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')

# test http gateway before the filter
curl -v -s -I "http://$INGRESS_HOST:$INGRESS_PORT/headers"

# install prime filter
kaf release/istio/filter

# test http gateway with the filter
curl  -v -s -I "http://$INGRESS_HOST:$INGRESS_PORT/headers"
