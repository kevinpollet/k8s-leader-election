CLUSTER_NAME=kle
CLUSTER_IMAGE=rancher/k3s:v1.22.3-k3s1

.PHONY: help
help:  ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: image
image: ## Build the Docker image.
	docker build -t leader-election .

.PHONY: deploy
deploy: ## Deploy the leader-election example to the current k8s cluster.
	kubectl apply -f deploy/

.PHONY: cluster
cluster: image ## Start a k3d cluster.
	k3d cluster create $(CLUSTER_NAME) --agents 2 --image $(CLUSTER_IMAGE) --k3s-server-arg '--no-deploy=traefik'
	k3d image import leader-election -c kle

.PHONY: cluster-delete
cluster-delete: ## Delete the k3d cluster.
	k3d cluster delete $(CLUSTER_NAME)
