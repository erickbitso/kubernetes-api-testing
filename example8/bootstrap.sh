kind create cluster --name kind

kubectl create namespace crossplane-system
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update
helm install crossplane --namespace crossplane-system crossplane-stable/crossplane

AWS_PROFILE=default && echo -e "[default]\naws_access_key_id = ${TEST_AWS_ACCESS_KEY_ID}\naws_secret_access_key = ${TEST_AWS_SECRET_ACCESS_KEY}" > creds.conf

kubectl create secret generic aws-creds -n crossplane-system --from-file=creds=$(pwd)/creds.conf
rm creds.conf
