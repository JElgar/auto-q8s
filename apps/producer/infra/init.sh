envsubst < cert.yml | kubectl apply -f -
envsubst < deployment.yml | kubectl apply -f -
envsubst < gateway.yml | kubectl apply -f -
envsubst < service.yml | kubectl apply -f -
envsubst < virtual_service.yml | kubectl apply -f -
