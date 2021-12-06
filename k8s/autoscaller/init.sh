envsubst < secret.yml | kubectl apply -f -
kubectl apply -f autoscaler.yml
