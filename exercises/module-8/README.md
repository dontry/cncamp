# httpserver - A simple HTTP server

### Deploy the httpserver pod and service

```sh
# wait for around 30 seconds for the pod to be ready
kubectl create -f httpserver-deploy.yaml
kubectl create -f httpserver-service.yaml
```

### Deploy fluentbit to collect logs from httpserver pods

```sh
kubectl create -f fluentbit.yaml
```

### View logs in fluentbit pod

```sh
kubectl logs -f fluentbit-xxxxx
```

### Deploy the ingress and ingress controller

```sh
# deploy k8s nginx ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.0/deploy/static/provider/cloud/deploy.yaml
```

### Generate key-cert

```sh
# use openssl (version >= 1.1.1f) on Linux, e.g. Ubuntu 20.04
# don't run on macOS, which is using LibreSSL
# instead, you can `brew install openssl` on macOS
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=httpserver.com/O=httpserver" -addext "subjectAltName = DNS:httpserver.com"
```

### Create secret

```sh
kubectl create secret tls httpserver-tls --cert=./tls.crt --key=./tls.key
```

### Deploy the ingress

```sh
kubectl create -f ingress.yaml
```

### Test the result

```sh
# get the ingress controller IP
kubectl get svc -n ingress-nginx

# test the ingress
curl -H "Host: httpserver.com" https://{INGRESS_CONTROLLER_IP} -v -k
```