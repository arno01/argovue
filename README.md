# kubevue

## Development

In skaffold folder rename `k8s-deploy.yaml.example` to `k8s-deploy.yaml` and provide missing values (OIDC), then:

```sh
skaffold dev --port-forward
```

You should see something like:

```sh
WARN[0000] port 50051 for gRPC server already in use: using 50052 instead
WARN[0000] port 50052 for gRPC HTTP server already in use: using 50053 instead
Listing files to watch...
 - kubevue
Generating tags...
 - kubevue -> kubevue:f9f2615-dirty
Checking cache...
 - kubevue: Found. Tagging
Tags used in deployment:
 - kubevue -> kubevue:2521378b4b22254b1fa4548a3d81af610bf049c24c557f137bac6ac6edbc5a14
   local images can't be referenced by digest. They are tagged and referenced by a unique ID instead
Starting deploy...
 - pod/kubevue created
Port forwarding pod/kubevue in namespace kubevue, remote port 8080 -> local port 8080
Watching for changes...
[kubevue main] time="2019-12-08T15:14:03Z" level=debug msg="Starting message broker"
[kubevue main] time="2019-12-08T15:14:04Z" level=debug msg="Starting kubernetes watcher: pods/v1/"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="Serving :8080, static folder:ui/dist"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add kube-apiserver-docker-desktop@kube-system uid:a332935f-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add compose-api-57ff65b8c7-tzjg5@docker uid:a3d6e4a7-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add compose-6c67d745f6-nsgzv@docker uid:a3e040ff-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add kubevue@kubevue uid:5da9a2bb-19cd-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add kube-proxy-hzvx2@kube-system uid:77346e0a-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add coredns-6dcc67dcbc-crnfx@kube-system uid:776aa5f9-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add etcd-docker-desktop@kube-system uid:a0378c5d-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add coredns-6dcc67dcbc-m2gxs@kube-system uid:7767fb7c-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add kube-scheduler-docker-desktop@kube-system uid:958104b3-19aa-11ea-bb60-025000000001"
[kubevue main] time="2019-12-08T15:14:04Z" level=info msg="add kube-controller-manager-docker-desktop@kube-system uid:96155666-19aa-11ea-bb60-025000000001"
```

After successful deployment point your browser to `http://localhost:8080/ui`.

