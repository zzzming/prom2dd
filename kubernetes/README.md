# Deployment on Kubernetes
Each tenant on each cluster should have a dedicated deployment. The tenant JWT must be specified in the Kubernete secret "{cluster}-{tenant}"

All pods are deployed under the DP cluster's kubernetes namespace `koddi`.

# Add secrets

kubectl create secret generic dd-app-key --from-literal=

kubectl create secret generic dd-api-key --from-literal=DD_API_KEY=

kubectl create secret generic dev-koddi-k1-azure-jwt --from-literal=PROMETHEUS_JWT_HEADER=eyJh...