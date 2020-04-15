go run reader.go instance.go parser.go pubsublistener.go worker.go -listen

Set environment variable
export GOOGLE_APPLICATION_CREDENTIALS=

export http_proxy without http://


export PROJECT_ID=[PROJECT_ID]

docker build -t gcr.io/${PROJECT_ID}/hello-app:v1 .
docker images
gcloud auth configure-docker

docker push gcr.io/${PROJECT_ID}/hello-app:v1

docker run --volume=/Users/andrepassos/projects/gocloudy/credentials/service-account.json:/secrets/service-account.json --env=GOOGLE_APPLICATION_CREDENTIALS=/secrets/service-account.json mono:latest

docker run -ti -v=/Users/andrepassos/.config/gcloud/:/root/.config/gcloud -p 8080:8080 go-docker

create secret from file:
kubectl create secret generic my-app-sa-key --from-file service-account.json

  PROJECT_ID=eleanor-270008
  534  gcloud config set compute/zone 
  535  gcloud config set compute/zone us-west1
  536  gcloud container clusters create hello-cluster --num-nodes=1
  537  gcloud compute instances list
  538  gcloud container clusters get-credentials hello-cluster
  540  kubectl create secret generic my-app-sa-key --from-file service-account.json
  541  kubectl get secret
  542  kubectl apply -f my-app.yaml 
  546  docker pull gcr.io/eleanor-270008/hello-app:v1
  547  kubectl apply -f my-app.yaml 