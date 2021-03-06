## Connection Secret Example

Small Golang app that executes `SET` and `GET` against Redis with random data when you hit it on `/redis`.

This is built to show how to move connection strings and passwords to a Kubernetes Secret - in order to:

1. Move Kubernetes applications from one service to another quickly.
2. You don't need to update **every** single application Deploy - just a single Secret.
3. Rollback is a single Secret file.
4. We can rollout connection changes application by application if desired. Can start with a small application - if it fails - we can quickly revert just that application.

To make a change:

1. Update the Secret with the new connection string AND/OR password.
2. Deploy that Secret.
3. Perform a `kubectl rollout restart deployment <name>`

Drawbacks?

1. A single secret change can affect multiple applications.
2. If a secret is deployed, the affected applications aren't automatically restarted.

## NOTES:

To create a Kubernetes cluster for testing on GKE:

```bash
gcloud beta container --project "project-name" \
    clusters create "my-test-cluster" \
    --region "us-east2" \
    --no-enable-basic-auth \
    --cluster-version "1.20.10-gke.1600" \
    --release-channel "stable" \
    --machine-type "e2-medium" \
    --image-type "COS_CONTAINERD" \
    --disk-type "pd-standard" \
    --disk-size "100" \
    --metadata disable-legacy-endpoints=true \
    --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" \
    --max-pods-per-node "110" \
    --num-nodes "3" \
    --logging=SYSTEM,WORKLOAD \
    --monitoring=SYSTEM \
    --enable-ip-alias \
    --network "projects/sre-scratchpad/global/networks/default" \
    --subnetwork "projects/project-name/regions/us-east2/subnetworks/default" \
    -no-enable-intra-node-visibility \
    --default-max-pods-per-node "110" \
    --no-enable-master-authorized-networks \
    --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver \
    --enable-autoupgrade \
    --enable-autorepair \
    --max-surge-upgrade 1 \
    --max-unavailable-upgrade 0 \
    --enable-shielded-nodes
```

To install test Redis on Kubernetes:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis \
    --set auth.password=$PASSWORD \
    --set replica.replicaCount=0 \
    --set tls.enabled=true \
    --set tls.authClients=false \
    --set tls.autoGenerated=true \
    --set metrics.enabled=true
```

To install Prometheus and Grafana

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack
```

TODO:

- [x] Test with Redis Password
- [x] Test with Redis TLS Connection
- [x] Kubernetes Secret
- [x] Kubernetes Deploy - healthz - annotations for Prometheus.
- [x] Kubernetes Service
- [x] Install Prometheus
- [x] Install Grafana
- [ ] Using Kustomize so we can deploy multiple copies of renamed services.
- [ ] Pods Monitored with Prometheus
- [ ] Redis Monitored with Prometheus
- [ ] Redis Cluster from Redis Labs
- [ ] Metrics from Redis Labs
- [ ] Demonstrate move from one cluster to another with metrics.
