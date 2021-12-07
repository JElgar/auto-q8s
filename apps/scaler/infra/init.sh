export SSH_PRIVATE_KEY=$(cat ~/.ssh/id_rsa)
envsubst < cert.yml | kubectl apply -f -
envsubst < deployment.yml | kubectl apply -f -
envsubst < gateway.yml | kubectl apply -f -
envsubst < service.yml | kubectl apply -f -
