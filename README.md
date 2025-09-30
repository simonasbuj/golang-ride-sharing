# golang-ride-sharing

bolt/uber wannabe using golang and event driven architecture
Please ignore "web" 

# tech and skills learned/improved in this project
- kubernetes with minikube
- tilt local kubernetes development tool
- rabbitmq
- websockets
- refactoring
- deploying to gcp Artifact Registry

# cheat sheat

### Make terminal user docker from minikube:
```
eval $(minikube docker-env)
eval $(minikube docker-env --unset)
```

### deploy to gcp
Get REGION and PROJECT_ID values from google cloud webpage.
Mine for example was REGION=europe-west1, PROJECT_ID=sandbox-big-query-437119

### Authenticate gcloud cli
```
gcloud auth login
gcloud auth configure-docker {REGION}-docker.pkg.dev
```

### Build images locally

```
docker build -t {REGION}-docker.pkg.dev/{PROJECT_ID}/golang-ride-sharing/api-gateway:latest --platform linux/amd64 -f infra/production/docker/api-gateway.Dockerfile .

docker build -t {REGION}-docker.pkg.dev/{PROJECT_ID}/golang-ride-sharing/trip-service:latest --platform linux/amd64 -f infra/production/docker/trip-service.Dockerfile .

```

### Push images to gcp Artifact Registry
```
docker push {REGION}-docker.pkg.dev/{PROJECT_ID}/golang-ride-sharing/api-gateway:latest

docker push {REGION}-docker.pkg.dev/{PROJECT_ID}/golang-ride-sharing/trip-service:latest
``` 

### Connect to gcp kubernetes
```
gcloud container clusters get-credentials golang-ride-sharing --region {REGION} --project {PROJECT_ID}
```

### Apply kubernetes files
```
kubectl apply -f infra/production/k8s/app-config.yaml

kubectl apply -f infra/production/k8s/api-gateway-deployment.yaml
kubectl apply -f infra/production/k8s/trip-service-deployment.yaml
```

Deploy all at once with:
```
kubectl apply -f infra/production/k8s/
```

### Call api
Get external-ip for `api-gateway` service by running `kubectl get services`.
Then you can make calls to apie using that ip for example `35.241.210.52:8081/trip/preview`
