#!/bin/bash
# -a is equivalent to allexport
set -a
source .env
echo "Launch GCP"
rm -rf ./terraform/vars.tf
envsubst < ./terraform/vars.tf.template > ./terraform/vars.tf
terraform init ./terraform
terraform apply ./terraform

echo "Deploy Elasticsearch and Kibana"
gcloud container clusters get-credentials $(terraform output cluster_name) --project ${GCP_PROJECT_ID} --zone $(terraform output cluster_zone)
# gcloud container clusters get-credentials ${APP}-cluster --project ${GCP_PROJECT_ID} --region ${GCP_REGION}
git clone https://github.com/ZedYeung/k8s-elk.git
cd ./k8s-elk
./run.sh
cd ..
rm -rf k8s-elk

while [[ -z "$ES_HOST"  ]];do
    ES_HOST=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" svc es-load-balancer -n elk)
    sleep 1
done
export ES_URL=http://${ES_HOST}:9200
echo "Elasticsearch URL: ${ES_URL}"

while [[ -z "$KIBANA_HOST"  ]];do
  KIBANA_HOST=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" svc kibana-load-balancer -n elk)
  sleep 1
done
export KIBANA_URL=http://${KIBANA_HOST}:9200
echo "KIBANA URL: ${KIBANA_URL}"

echo "Deploy backend and frontend"
rm -rf docker-compose.yml
envsubst < docker-compose.yml.template > docker-compose.yml
# kompose up docker-compose.yml
kompose up docker-compose.yml

echo "Deploy Ingress"
# https
kubectl create secret tls ${APP}-tls-cert --key ../acme/certificate_pem/${DOMAIN}_private_key.pem --cert ../acme/certificate_pem/${DOMAIN}.pem
# export certificate_body=$(cat ../acme/certificate_pem/${DOMAIN}_private_key.pem | base64 -w 0)
# export private_key=$(cat ../acme/certificate_pem/${DOMAIN}.pem | base64 -w 0)
# envsubst < ./secret.yml.tmpl > ./secret.yml
# kubectl apply -f secret.yml

# ingress-nginx
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml

envsubst < ingress.yml.tmpl > ingress.yml
kubectl apply -f ingress.yml
echo "Done"