kubectl create ns rabbits
kubectl apply -n rabbits -f rbac.yml
kubectl apply -n rabbits -f configmap.yml
kubectl apply -n rabbits -f secret.yml
kubectl apply -n rabbits -f statefullset.yml
