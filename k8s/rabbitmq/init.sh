kubectl apply -f queue.yml
envsubst < gateway.yml | kubectl apply -f -
