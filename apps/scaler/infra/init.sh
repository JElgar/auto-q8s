export SSH_PRIVATE_KEY=$(cat ~/.ssh/id_rsa)
envsubst < secret.yml | kubectl apply -f -
envsubst < deployment.yml | kubectl apply -f -
