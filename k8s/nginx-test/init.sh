envsubst < nginx-cert.yml | kubectl apply -f -
envsubst < nginx-deployment.yml | kubectl apply -f -
envsubst < nginx-gateway.yml | kubectl apply -f -
envsubst < nginx-service.yml | kubectl apply -f -
