# Rathole Operator Helm Chart

Helm chart will move to operator repository.

See: [Rathole Operator](https://github.com/crmin/rathole-operator)

## Install Helm Chart

First, add the repository to your helm client.
```
$ helm repo add rathole-operator https://rathole-operator.superclass.io
$ helm repo update
$ helm repo list
$ helm search repo rathole-operator
```

If repo added successfully, you can see the chart in the list.

```
NAME                             	CHART VERSION	APP VERSION	DESCRIPTION
rathole-operator/rathole-operator	0.1.0        	1.16.0     	A Helm chart for Kubernetes
```

Then, install the chart with the following command.

```
$ helm install rathole-operator rathole-operator/rathole-operator
```

or you can install with the values file.

```
$ helm install rathole-operator rathole-operator/rathole-operator -f values.yaml
```

## Values

```yaml
crdVersion: v1alpha1

serviceAccount:
  create: true
  name: rathole-operator

replicaCount: 1

image:
  repository: crmin/rathole-operator
  pullPolicy: IfNotPresent
  tag: "v1alpha1"

affinity: {}
nodeSelector: {}

resources:
  limits:
    cpu: 500m
    memory: 128Mi
  requests:
    cpu: 10m
    memory: 64Mi

namespace:
  create: true
  name: system

imagePullSecrets: []
livenessProbe:
  httpGet:
    path: /healthz
    port: 8081
  initialDelaySeconds: 15
  periodSeconds: 20
readinessProbe:
  httpGet:
    path: /readyz
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 10
```
