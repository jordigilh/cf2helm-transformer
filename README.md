# cf2helm-transformer
CF Application discovery to helm transformer

This is a sample helm chart that uses the values.yaml generated from a CF Application discovery structure and renders 3 simple manifests:
* Deployment
* Service
* Ingress

Running `helm template sample sample` will render a yaml list that contains the 3 Kubernetes manifests populated with the values in the
`values.yaml` file.

# Building
To build the binary, run the following command:
```shell
go build
```

# Generating the `values.yaml` file
It will create a new binary named `cf2helm-transformer`. Running the binary will generate the `values.yaml` in the `sample/` directory as well
as output the contents of the file generated.

```shell
$>./cf2helm-transformer
name: foo
space: default
labels:
  foo: bar
annotations:
  bar: foo
env:
  CF_ROOT: /root
  SERVER_PORT: "8080"
route:
- route: foo.default.cluster.io/
  protocol: http2
service:
- name: mysql
  parameters:
    DB_NAME: default
process:
- type: web
  command:
  - /usr/bin/echo hello world>index.html; /usr/local/bin/python3 -m http.server $SERVER_PORT
  memory: 128Mi
  healthCheck:
    endpoint: localhost:8080/
    type: http
  readinessCheck:
    endpoint: localhost:8080/
    type: http
stack: default
docker:
  image: python:3
instances: 1

```

# Rendering the Kubernetes templates
Run the following command to render the Kubernetes manifests, based on the contents of the `values.yaml`

```
helm template sample sample
```

Example:
```
$>helm template sample sample
---
# Source: sample/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  namespace: default
  name: foo
  labels:
    app: foo
    foo: bar
  annotations:
    bar: foo
    stack: default
spec:
  selector:
    app: foo
  ports:
    - targetPort: 8080
      protocol: TCP
      port: 80
      name: web
---
# Source: sample/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: foo
  labels:
    app: foo
    foo: bar
  annotations:
    bar: foo
    stack: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: foo
  template:
    metadata:
      namespace: default
      name: foo
      labels:
        app: foo
        foo: bar
      annotations:
        bar: foo
        stack: default
    spec:
      containers:
        - name: web
          command:
            - /bin/sh
          args:
            - -c
            - /usr/bin/echo hello world>index.html; /usr/local/bin/python3 -m http.server $SERVER_PORT
          env:
            - name: CF_ROOT
              value: "/root"
            - name: SERVER_PORT
              value: "8080"
            - name: SERVICE__FOO__DB_NAME
              value: "default"
          image: python:3
          ports:
            - name: web
              containerPort: 8080
          livenessProbe:
            httpGet:
              path: "/"
              port: 8080
          readinessProbe:
            httpGet:
              path: "/"
              port: 8080
          resources:
            requests:
              memory: 128Mi
              cpu: 100m
---
# Source: sample/templates/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: default
  name: foo
  labels:
    app: foo
    foo: bar
  annotations:
    bar: foo
    stack: default
spec:
  rules:
    - host: foo.default.cluster.io
      http:
        paths:
          - path: /
            pathType: Exact
            backend:
              service:
                name: foo
                port:
                  number: 8080
```
Note that to generate the templates, you will need to have the `helm` client installed.
