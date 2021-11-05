# k8s-leader-election

Leader election example leveraging the [leaderelection](https://github.com/kubernetes/client-go/tree/master/tools/leaderelection) pkg from the k8s [client](https://github.com/kubernetes/client-go#readme).

## Usage

1. Create a k3d cluster: `make cluster`
2. Deploy the leader-election app: `make deploy`
3. Get the current lease holder:

```kubectl get lease leader-election -o=jsonpath='{.spec.holderIdentity}'```

5. Kill the current lease holder:

```kubectl delete pod $(kubectl get lease leader-election -o=jsonpath='{.spec.holderIdentity}')```

7. After waiting the `leaseDuration` period (`60s`) a new leader will be elected.

## License

[MIT](./LICENSE.md)
